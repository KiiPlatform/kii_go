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
func (gb GroupBucket) Path() string {
	return "/groups/" + gb.GroupID + "/buckets/" + gb.BucketName
}

// ThingBucket represents thing bucket
type ThingBucket struct {
	BucketName string
	ThingID    string
}

// Path generates path string of user bucket
func (tb ThingBucket) Path() string {
	return "/things/" + tb.ThingID + "/buckets/" + tb.BucketName
}
