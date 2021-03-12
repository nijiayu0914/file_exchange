package services

import (
	"errors"
	"file_exchange/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type OssOperator struct {
	Endpoint        		string
	AccessKeyId     		string
	AccessKeySecret 		string
	OSSBucketName 			string
	OSSRegionId 			string
	OSSRamAccessKeyID 		string
	OSSRamAccessKeySecret 	string
	OSSRoleArn 				string
	OSSRoleSessionName		string
	Client          		*oss.Client
	Bucket 					*oss.Bucket
}

func (os *OssOperator) GetClient () (err error){
	os.Client, err = oss.New(os.Endpoint, os.AccessKeyId, os.AccessKeySecret)
	if err != nil {
		return err
	}
	os.Bucket, err = os.Client.Bucket(os.OSSBucketName)
	if err != nil {
		return nil
	}
	 return nil
}

func (os *OssOperator) CreateSTS () (response *sts.AssumeRoleResponse, err error){
	client, err := sts.NewClientWithAccessKey(
		os.OSSRegionId, os.OSSRamAccessKeyID, os.OSSRamAccessKeySecret)
	if err != nil{
		return nil, err
	}
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = os.OSSRoleArn
	request.RoleSessionName = os.OSSRoleSessionName
	response, err = client.AssumeRole(request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (os *OssOperator) CreateStringFile (fileUuid string,
	childFileName string, content string, category string) (err error) {
	options := []oss.Option{
		oss.ObjectStorageClass(oss.StorageStandard),
		oss.ContentDisposition(childFileName),
		oss.ObjectACL(oss.ACLPrivate),
		oss.Meta("category", category),
		oss.Meta("origin-name", childFileName),
		oss.Meta("rename", childFileName),
	}

	err = os.Bucket.PutObject(
		utils.StringJoin("/", fileUuid, childFileName),
		strings.NewReader(content), options...)
	if err != nil {
		return err
	}
	return nil
}

func (os *OssOperator) GetObjectMetaData (fileName string, versionId string) (http.Header, error) {
	props, err := os.Bucket.GetObjectDetailedMeta(fileName, oss.VersionId(versionId))
	if err != nil {
		return nil, err
	}
	return props, nil
}

func (os *OssOperator) RenameObject (objectName string, newName string) (err error){
	props, err := os.Bucket.GetObjectDetailedMeta(objectName)
	if err != nil{
		return err
	}
	err = os.Bucket.SetObjectMeta(
		objectName,
		oss.ContentDisposition(props.Get("Content-Disposition")),
		oss.Meta("rename", newName),
		oss.Meta("category", props.Get("x-oss-meta-category")),
		oss.Meta("origin-name", props.Get("x-oss-meta-origin-name")),
	)
	if err != nil {
		return err
	}
	return nil
}

func (os *OssOperator) ListFiles (fileUuid string,
	path string, delimiter string)  (objectsContainer []utils.ObjectInfoCollection,
	dirsContainer []utils.DirInfoCollection,
	err error) {
	continueToken := ""
	for {
		lsRes, err := os.Bucket.ListObjectsV2(
			oss.Prefix(utils.StringJoin("/", fileUuid, path)),
			oss.ContinuationToken(continueToken),
			oss.Delimiter(delimiter),
			oss.MaxKeys(1000))
		if err != nil {
			return objectsContainer, dirsContainer, err
		}

		for _, object := range lsRes.Objects {
			info := utils.ObjectInfoCollection{}
			metaData, err := os.Bucket.GetObjectDetailedMeta(object.Key)
			if err != nil{
				info.Meta = http.Header{}
			}else{
				info.Meta = metaData
			}
			info.Basic = object
			objectsContainer = append(objectsContainer, info)
		}
		for _, dirName := range lsRes.CommonPrefixes {
			info := utils.DirInfoCollection{}
			metaData, err := os.Bucket.GetObjectDetailedMeta(dirName)
			if err != nil{
				info.Meta = http.Header{}
			}else{
				info.Meta = metaData
			}
			info.Basic = dirName
			dirsContainer = append(dirsContainer, info)
		}

		if lsRes.IsTruncated {
			continueToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
	return objectsContainer, dirsContainer, err
}

func (os *OssOperator) IsExist (fileUuid string,
	fileName string) (isExist bool, err error){
	isExist, err = os.Bucket.IsObjectExist(
		utils.StringJoin("/", fileUuid, fileName))
	if err != nil {
		return false, err
	}
	return isExist, nil
}

func (os *OssOperator) DeleteFile (fileUuid string,
	fileName string) (err error) {
	err = os.Bucket.DeleteObject(utils.StringJoin("/", fileUuid, fileName))
	if err != nil {
		return err
	}
	return nil
}

func (os *OssOperator) DeleteChildFile (fileUuid string,
	path string) (err error){
	objectsContainer, _, err := os.ListFiles(fileUuid, path, "")
	for _, object := range objectsContainer{
		err = os.Bucket.DeleteObject(object.Basic.Key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (os *OssOperator) DeleteFiles (fileUuid string,
	fileNames []string) (deleteMarket []oss.DeletedKeyInfo, err error) {
	var delObjects []oss.DeleteObject
	for _, fileName := range fileNames{
		delObjects = append(delObjects,
			oss.DeleteObject{Key: utils.StringJoin("/", fileUuid, fileName)})
	}
	res, err := os.Bucket.DeleteObjectVersions(delObjects)
	if err != nil {
		return nil, err
	}
	deleteMarket = res.DeletedObjectsDetail
	return deleteMarket, nil
}

func(os *OssOperator) ListFileVersion (fileUuid string,
	path string) (objects []map[string]interface{}, err error){
	delimiter := oss.Delimiter("/")
	keyMarker := oss.KeyMarker("")
	versionIdMarker := oss.VersionIdMarker("")
	prefix := oss.Prefix(utils.StringJoin("/", fileUuid, path))
	for {
		lor, err := os.Bucket.ListObjectVersions(prefix, keyMarker, delimiter, versionIdMarker)
		if err != nil {
			return nil, err
		}
		for _, obj  := range lor.ObjectVersions {
			object := map[string]interface{}{
				"versionId": obj.VersionId,
				"key": obj.Key,
				"isLatest": obj.IsLatest,
				"etag": obj.ETag,
				"size": obj.Size,
				"lastModified": obj.LastModified,
			}
			objects = append(objects, object)
		}
		if lor.IsTruncated {
			keyMarker = oss.KeyMarker(lor.NextKeyMarker)
			versionIdMarker = oss.VersionIdMarker(lor.NextVersionIdMarker)
		}else{
			break
		}
	}
	return objects, nil
}

func(os *OssOperator) ListDeleteMarkers (fileUuid string, path string, delimiter string) (
	markers []map[string]interface{}, err error) {
	dmt := oss.Delimiter(delimiter)
	keyMarker := oss.KeyMarker("")
	versionIdMarker := oss.VersionIdMarker("")
	prefix := oss.Prefix(utils.StringJoin("/", fileUuid, path))
	for {
		lor, err := os.Bucket.ListObjectVersions(prefix, keyMarker, dmt, versionIdMarker)
		if err != nil {
			return nil, err
		}
		for _, mk  := range lor.ObjectDeleteMarkers {
			marker := map[string]interface{}{
				"versionId": mk.VersionId,
				"key": mk.Key,
				"isLatest": mk.IsLatest,
				"lastModified": mk.LastModified,
			}
			markers = append(markers, marker)
		}
		if lor.IsTruncated {
			keyMarker = oss.KeyMarker(lor.NextKeyMarker)
			versionIdMarker = oss.VersionIdMarker(lor.NextVersionIdMarker)
		}else{
			break
		}
	}
	return markers, nil
}

func (os *OssOperator) DeleteHistoryFile(fileUuid string,
	path string, versionId string) (err error){
	key := utils.StringJoin("/", fileUuid, path)
	var delObjects []oss.DeleteObject
	delObjects = append(
		delObjects,
		oss.DeleteObject{Key: key, VersionId: versionId})
	_, err =os.Bucket.DeleteObjectVersions(delObjects)
	if err != nil {
		return err
	}
	return nil
}

func(os *OssOperator) DeleteFileForever(
	fileUuid string, fileName string) (size float64, err error){
	markers, err := os.ListDeleteMarkers(fileUuid, fileName, "/")
	if err != nil{
		return 0, err
	}
	objects, err := os.ListFileVersion(fileUuid, fileName)
	if err != nil{
		return 0, err
	}
	var delObjects []oss.DeleteObject
	for _, deleteObj := range markers{
		delObjects = append(
			delObjects,
			oss.DeleteObject{Key: deleteObj["key"].(string),
				VersionId: deleteObj["versionId"].(string)})
	}
	for _, deleteObj := range objects{
		delObjects = append(
			delObjects,
			oss.DeleteObject{Key: deleteObj["key"].(string),
				VersionId: deleteObj["versionId"].(string)})
		delSize, _ := os.ReadFileCapacity(
			deleteObj["key"].(string), deleteObj["versionId"].(string))
		size += delSize
	}
	_, err = os.Bucket.DeleteObjectVersions(delObjects)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func(os *OssOperator) DeleteFilesForever(fileUuid string,
	fileNames []string) (size float64, err error) {
	for _, fileName := range fileNames{
		Delsize, err := os.DeleteFileForever(fileUuid, fileName)
		size += Delsize
		if err != nil {
			return 0, err
		}
	}
	return size, nil
}

func(os *OssOperator) Copy(originFile string,
	destFile string, versionId string) (addSize float64, err error) {
	originFileInfo, err := os.GetObjectMetaData(originFile, versionId)
	if err != nil{
		return 0, err
	}
	size, err := strconv.ParseFloat(originFileInfo.Get("Content-Length"), 64)
	if err != nil{
		return 0, err
	}
	mbSize := size / 1024 / 1024
	if mbSize <= 900{
		_, err = os.Bucket.CopyObject(
			originFile,
			destFile,
			oss.VersionId(versionId))
		if err != nil {
			return 0, err
		}
	}else{
		chunkSize := int(math.Floor(mbSize / 500))
		chunks, err := oss.SplitFileByPartNum(originFile, chunkSize)
		if err != nil{
			return 0, err
		}
		var parts []oss.UploadPart
		imu, err := os.Bucket.InitiateMultipartUpload(destFile)
		for _, chunk := range chunks {
			part, err := os.Bucket.UploadPartCopy(
				imu, os.OSSBucketName, originFile,
				chunk.Offset, chunk.Size, chunk.Number,
				oss.VersionId(versionId))
			if err != nil {
				return 0, err
			}
			parts = append(parts, part)
		}
		_, err = os.Bucket.CompleteMultipartUpload(imu, parts)
		if err != nil {
			return 0, err
		}
	}
	addSize, _ = os.ReadFileCapacity(originFile, versionId)
	return addSize, nil
}

func(os *OssOperator) MultipleCopy(copyList []utils.RequestCopy)(
	failure []utils.RequestCopy, size float64, err error) {
	for _, cp := range copyList{
		addSize, err := os.Copy(cp.OriginFile, cp.DestFile, cp.VersionId)
		size += addSize
		if err != nil{
			failure = append(failure, cp)
		}
	}
	if len(failure) == 0{
		return failure, size, nil
	}else{
		return failure, 0, errors.New("复制失败")
	}
}

func(os *OssOperator) RestoreFile(fileUuid string, path string) (err error){
	markers, err := os.ListDeleteMarkers(fileUuid, path, "/")
	if err != nil{
		return nil
	}
	var delObjects []oss.DeleteObject
	for _, deleteObj := range markers{
		delObjects = append(
			delObjects,
			oss.DeleteObject{Key: deleteObj["key"].(string),
				VersionId: deleteObj["versionId"].(string)})
	}
	_, err = os.Bucket.DeleteObjectVersions(delObjects)
	if err != nil {
		return err
	}
	return nil
}

func(os *OssOperator) ReadFileCapacity(fileName string,
	versionId string) (size float64, err error){
	delimiter := oss.Delimiter("/")
	keyMarker := oss.KeyMarker("")
	versionIdMarker := oss.VersionIdMarker("")
	prefix := oss.Prefix(fileName)
	for {
		lor, err := os.Bucket.ListObjectVersions(prefix, keyMarker, delimiter, versionIdMarker)
		if err != nil {
			return 0, err
		}
		for _, obj  := range lor.ObjectVersions {
			if obj.VersionId == versionId{
				return float64(obj.Size) / 1024 / 1024, nil
			}
		}
		if lor.IsTruncated {
			keyMarker = oss.KeyMarker(lor.NextKeyMarker)
			versionIdMarker = oss.VersionIdMarker(lor.NextVersionIdMarker)
		}else{
			break
		}
	}
	return 0, errors.New("未找到文件")
}

func(os *OssOperator) ReadFilesCapacity(fileUuid string,
	files []utils.RequestReadFileSize) (size float64, err error){
	for _, obj := range files{
		objectName := utils.StringJoin("/", fileUuid, obj.FileName)
		objSize, err := os.ReadFileCapacity(objectName, obj.VersionId)
		if err != nil{
			return size, err
		}
		size += objSize
	}
	return size, nil
}

func(os *OssOperator) ReadAllFilesCapacity(fileUuid string) (size float64, err error){
	objects, _, _ := os.ListFiles(fileUuid, "", "")
	for _, obj := range objects{
		fileName := strings.Replace(obj.Basic.Key, fileUuid + "/", "", -1)
		objectsVs, err := os.ListFileVersion(fileUuid, fileName)
		if err != nil{
			return 0, err
		}
		for _, ovs := range objectsVs{
			objSize, err := os.ReadFileCapacity(obj.Basic.Key, ovs["versionId"].(string))
			if err != nil{
				return 0, err
			}
			size += objSize
		}
	}
	return size, nil
}

func(os *OssOperator) DownloadFile(fileUuid string,
	fileName string) (url string, err error){
	objectName := utils.StringJoin("/", fileUuid, fileName)
	signedURL, err := os.Bucket.SignURL(objectName, oss.HTTPGet, 60)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}
