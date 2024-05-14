/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package wasabicloud

type WasabiStorageHandler struct{ s3Conf }

func NewWasabiHandler() *WasabiStorageHandler {
	return &WasabiStorageHandler{
		s3Conf{
			bucketName: "avx",
			region:     "ap-southeast-1",
			endpoint:   "s3.ap-southeast-1.wasabisys.com",
			accessKey:  "5XNJOEBTPNNOWU9QBOMK",
			secretKey:  "mHeOikDwKF5WMeLvHKeyVsVeHqvwKiFcrKOr5MGR",
		},
	}
}
