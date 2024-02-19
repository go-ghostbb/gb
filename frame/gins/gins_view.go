package gins

import (
	"context"
	"fmt"
	"ghostbb.io/gb/internal/consts"
	"ghostbb.io/gb/internal/instance"
	"ghostbb.io/gb/internal/intlog"
	gbview "ghostbb.io/gb/os/gb_view"
	gbutil "ghostbb.io/gb/util/gb_util"
)

// View returns an instance of View with default settings.
// The parameter `name` is the name for the instance.
// Note that it panics if any error occurs duration instance creating.
func View(name ...string) *gbview.View {
	instanceName := gbview.DefaultName
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameViewer, instanceName)
	return instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		return getViewInstance(instanceName)
	}).(*gbview.View)
}

func getViewInstance(name ...string) *gbview.View {
	var (
		err          error
		ctx          = context.Background()
		instanceName = gbview.DefaultName
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	view := gbview.Instance(instanceName)
	if Config().Available(ctx) {
		var (
			configMap      map[string]interface{}
			configNodeName = consts.ConfigNodeNameViewer
		)
		if configMap, err = Config().Data(ctx); err != nil {
			intlog.Errorf(ctx, `retrieve config data map failed: %+v`, err)
		}
		if len(configMap) > 0 {
			if v, _ := gbutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameViewer); v != "" {
				configNodeName = v
			}
		}
		configMap = Config().MustGet(ctx, fmt.Sprintf(`%s.%s`, configNodeName, instanceName)).Map()
		if len(configMap) == 0 {
			configMap = Config().MustGet(ctx, configNodeName).Map()
		}
		if len(configMap) > 0 {
			if err = view.SetConfigWithMap(configMap); err != nil {
				panic(err)
			}
		}
	}
	return view
}
