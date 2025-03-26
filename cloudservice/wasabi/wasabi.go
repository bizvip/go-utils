package wasabi

// var Example = conf{
// 	bucketName: "avx",
// 	region:     "ap-southeast-1",
// 	endpoint:   "s3.ap-southeast-1.wasabisys.com",
// 	accessKey:  "5XNJOEBTPNNOWU9QBOMK",
// 	secretKey:  "mHeOikDwKF5WMeLvHKeyVsVeHqvwKiFcrKOr5MGR",
// }

// conf 结构体包含 Wasabi S3 的配置信息
type conf struct {
	bucketName string
	region     string
	endpoint   string
	accessKey  string
	secretKey  string
}

// StorageHandler 结构体，用于操作 Wasabi S3 存储
type StorageHandler struct {
	s3Conf conf
}

// NewWasabiHandler 创建一个新的 StorageHandler 实例，并初始化配置信息
func NewWasabiHandler(bucketName, region, endpoint, accessKey, secretKey string) *StorageHandler {
	return &StorageHandler{
		s3Conf: conf{
			bucketName: bucketName,
			region:     region,
			endpoint:   endpoint,
			accessKey:  accessKey,
			secretKey:  secretKey,
		},
	}
}
