package kii

// Bucket represents bucket
type Bucket interface {
	Path() string
}

// AppBucket represents app bucket
type AppBucket struct {
	BucketName string
}

// Path generates path string of app bucket
func (ab AppBucket) Path() string {
	return "/buckets/" + ab.BucketName
}

// UserBucket represents user bucket
type UserBucket struct {
	BucketName string
	UserID     string
}

// Path generates path string of user bucket
func (ub UserBucket) Path() string {
	return "/users/" + ub.UserID + "/buckets/" + ub.BucketName
}

// GroupBucket represents group bucket
type GroupBucket struct {
	BucketName string
	GroupID    string
}

// Path generates path string of group bucket
func (ub GroupBucket) Path() string {
	return "/groups/" + ub.GroupID + "/buckets/" + ub.BucketName
}

// ThingBucket represents thing bucket
type ThingBucket struct {
	BucketName string
	ThingID    string
}

// Path generates path string of user bucket
func (ub ThingBucket) Path() string {
	return "/things/" + ub.ThingID + "/buckets/" + ub.BucketName
}
