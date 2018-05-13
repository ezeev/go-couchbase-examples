package bestbuy

import gocb "gopkg.in/couchbase/gocb.v1"

func OpenCluster(address string, username string, password string) (*gocb.Cluster, error) {
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

func OpenBucket(name string, cluster *gocb.Cluster) (*gocb.Bucket, error) {
	bucket, err := cluster.OpenBucket(name, "")
	return bucket, err
}

func MustOpenCluster(address string, username string, password string) *gocb.Cluster {
	cluster, err := OpenCluster(address, username, password)
	if err != nil {
		panic(err)
	}
	return cluster
}

func MustOpenBucket(name string, cluster *gocb.Cluster) *gocb.Bucket {
	bucket, err := OpenBucket(name, cluster)
	if err != nil {
		panic(err)
	}
	return bucket
}
