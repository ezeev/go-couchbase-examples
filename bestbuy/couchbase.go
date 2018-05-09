package bestbuy

import gocb "gopkg.in/couchbase/gocb.v1"

func CbConnect(address string, username string, password string) (*gocb.Cluster, error) {
	cluster, err := gocb.Connect(address)
	if err != nil {
		return nil, err
	}
	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func CbOpenBucket(name string, cluster *gocb.Cluster) (*gocb.Bucket, error) {
	bucket, err := cluster.OpenBucket(name, "")
	return bucket, err
}
