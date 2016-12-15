package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/smartystreets/go-aws-auth"
)

const ObjectStorageURI = "https://%s.os.cloud.nifty.com"

// Plugin defines the S3 plugin parameters.
type Plugin struct {
	Key         string
	Secret      string
	Bucket      string
	Region      string
	Access      string
	Source      string
	Target      string
	StripPrefix string
	Exclude     []string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	// initialize storage client
	credentials := awsauth.Credentials{
		AccessKeyID:     p.Key,
		SecretAccessKey: p.Secret,
	}
	client := NewClient(fmt.Sprintf(ObjectStorageURI, p.Region), credentials)

	// normalize the target URL
	if strings.HasPrefix(p.Target, "/") {
		p.Target = p.Target[1:]
	}

	// find the bucket
	log.WithFields(log.Fields{
		"region": p.Region,
		"bucket": p.Bucket,
	}).Info("Attempting to upload")

	matches, err := matches(p.Source, p.Exclude)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not match files")
		return err
	}

	for _, match := range matches {

		stat, err := os.Stat(match)
		if err != nil {
			continue // should never happen
		}

		// skip directories
		if stat.IsDir() {
			continue
		}

		target := filepath.Join(p.Target, strings.TrimPrefix(match, p.StripPrefix))

		// log file for debug purposes.
		log.WithFields(log.Fields{
			"name":   match,
			"bucket": p.Bucket,
			"target": target,
		}).Info("Uploading file")

		// put object
		err = client.PutObject(p.Bucket, target, match, p.Access)
		if err != nil {
			log.WithFields(log.Fields{
				"name":   match,
				"bucket": p.Bucket,
				"target": target,
				"error":  err,
			}).Error("Could not upload file")
			return err
		}
	}
	return nil
}

// matches is a helper function that returns a list of all files matching the
// included Glob pattern, while excluding all files that matche the exclusion
// Glob pattners.
func matches(include string, exclude []string) ([]string, error) {
	matches, err := zglob.Glob(include)
	if err != nil {
		return nil, err
	}
	if len(exclude) == 0 {
		return matches, nil
	}

	// find all files that are excluded and load into a map. we can verify
	// each file in the list is not a member of the exclusion list.
	excludem := map[string]bool{}
	for _, pattern := range exclude {
		excludes, err := zglob.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, match := range excludes {
			excludem[match] = true
		}
	}

	var included []string
	for _, include := range matches {
		_, ok := excludem[include]
		if ok {
			continue
		}
		included = append(included, include)
	}
	return included, nil
}
