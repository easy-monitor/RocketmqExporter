package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"RocketmqExporter/config"
	"RocketmqExporter/constant"
	"RocketmqExporter/utils"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	RocketmqConsoleIPAndPort string
	JSESSIONID               string
	last_target              string
)

func loadConfig() (*config.Conf, error) {
	path, _ := os.Getwd()
	path = filepath.Join(path, "conf/conf.yml")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("read conf.yml fail:" + path)
	}
	conf := new(config.Conf)
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, errors.New("unmarshal conf.yml fail" + err.Error())
	}
	return conf, nil
}

func main() {
	RocketmqConsoleIPAndPort = "localhost:80"
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)

	kingpin.Version(version.Print("rocketmq_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting rocketmq_exporter", "version", version.Info)
	level.Info(logger).Log("msg", "Build contenxt", version.BuildContext())

	conf, err := loadConfig()
	if err != nil {
		fmt.Println("loadConfig fail:" + err.Error())
		level.Error(logger).Log("err", "loadConfig fail:"+err.Error())
	}

	metricsPath := constant.GetMetricsPath()
	RocketmqConsoleIPAndPort = "localhost:8080"
	listenAddress := ":" + conf.Port
	constant.SetIgnoredTopicArray(conf.IgnoredTopics)
	metricsPrefix := constant.GetMetricsPrefix()

	level.Info(logger).Log("msg", "fmt.metricsPath:"+metricsPath)

	exporter := DeclareExporter(metricsPrefix)
	prometheus.MustRegister(exporter)

	http.Handle("/scrape", ScrapeHandlerFor(metricsPrefix, conf))

	http.Handle(metricsPath, promhttp.Handler())
	fmt.Println(http.ListenAndServe(listenAddress, nil))
}

func ScrapeHandlerFor(metricsPrefix string, conf *config.Conf) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Query()
		target := uri.Get("target")
		module := uri.Get("module")
		if target == "" || module == "" {
			buf := "args error"
			_, _ = w.Write([]byte(buf))
			return
		}
		//登出上一个登陆
		if last_target != "" {
			logoutUrl := "http://" + last_target + "/login/logout.do"
			utils.HttpUrl("POST", logoutUrl, JSESSIONID)
		}
		//重新登陆
		modules := conf.Module
		var username = ""
		var password = ""
		for _, v := range modules {
			if v.Name == module {
				username = v.Username
				password = v.Password
				break
			}
		}
		if username == "" || password == "" {
			buf := "module error, not found module in conf.yml"
			_, _ = w.Write([]byte(buf))
			return
		}
		loginUrl := "http://" + target + "/login/login.do?username=" + username + "&password=" + password
		rep := utils.HttpUrl("POST", loginUrl, "")
		loginResponse := config.LoginResponse{}
		err := json.Unmarshal(rep, &loginResponse)
		if err != nil {
			buf := "login fail"
			_, _ = w.Write([]byte(buf))
			return
		}
		JSESSIONID = loginResponse.Data.SessionId
		RocketmqConsoleIPAndPort = target
		exporter := DeclareExporter(metricsPrefix)
		registry := prometheus.NewRegistry()
		_ = registry.Register(exporter)
		gatherers := prometheus.Gatherers{
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
		last_target = target
	})
}
