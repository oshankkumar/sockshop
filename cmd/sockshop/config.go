package main

import "flag"

type AppConfig struct {
	MySQLConnString string
	ImagePath       string
	Domain          string
}

func NewConfigFromFlags() AppConfig {
	var conf AppConfig
	flag.StringVar(&conf.MySQLConnString, "mysql-conn-str", "admin:password@tcp(mysql:3306)/socksdb", "MySQL connection string")
	flag.StringVar(&conf.ImagePath, "image-path", "assets/images", "Image path")
	flag.StringVar(&conf.Domain, "link-domain", "127.0.0.1:9090", "HATEAOS link domain")
	flag.Parse()
	return conf
}
