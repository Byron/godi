package server

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

// bower_components_angular_angular_csp_css reads file data from disk. It returns an error on failure.
func bower_components_angular_angular_csp_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/angular-csp.css",
		"bower_components/angular/angular-csp.css",
	)
}

// bower_components_angular_angular_js reads file data from disk. It returns an error on failure.
func bower_components_angular_angular_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/angular.js",
		"bower_components/angular/angular.js",
	)
}

// bower_components_angular_angular_min_js reads file data from disk. It returns an error on failure.
func bower_components_angular_angular_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/angular.min.js",
		"bower_components/angular/angular.min.js",
	)
}

// bower_components_angular_angular_min_js_gzip reads file data from disk. It returns an error on failure.
func bower_components_angular_angular_min_js_gzip() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/angular.min.js.gzip",
		"bower_components/angular/angular.min.js.gzip",
	)
}

// bower_components_angular_angular_min_js_map reads file data from disk. It returns an error on failure.
func bower_components_angular_angular_min_js_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/angular.min.js.map",
		"bower_components/angular/angular.min.js.map",
	)
}

// bower_components_angular_bower_json reads file data from disk. It returns an error on failure.
func bower_components_angular_bower_json() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/bower.json",
		"bower_components/angular/bower.json",
	)
}

// bower_components_angular_readme_md reads file data from disk. It returns an error on failure.
func bower_components_angular_readme_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular/README.md",
		"bower_components/angular/README.md",
	)
}

// bower_components_angular_loader_angular_loader_js reads file data from disk. It returns an error on failure.
func bower_components_angular_loader_angular_loader_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-loader/angular-loader.js",
		"bower_components/angular-loader/angular-loader.js",
	)
}

// bower_components_angular_loader_angular_loader_min_js reads file data from disk. It returns an error on failure.
func bower_components_angular_loader_angular_loader_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-loader/angular-loader.min.js",
		"bower_components/angular-loader/angular-loader.min.js",
	)
}

// bower_components_angular_loader_angular_loader_min_js_map reads file data from disk. It returns an error on failure.
func bower_components_angular_loader_angular_loader_min_js_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-loader/angular-loader.min.js.map",
		"bower_components/angular-loader/angular-loader.min.js.map",
	)
}

// bower_components_angular_loader_bower_json reads file data from disk. It returns an error on failure.
func bower_components_angular_loader_bower_json() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-loader/bower.json",
		"bower_components/angular-loader/bower.json",
	)
}

// bower_components_angular_loader_readme_md reads file data from disk. It returns an error on failure.
func bower_components_angular_loader_readme_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-loader/README.md",
		"bower_components/angular-loader/README.md",
	)
}

// bower_components_angular_mocks_angular_mocks_js reads file data from disk. It returns an error on failure.
func bower_components_angular_mocks_angular_mocks_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-mocks/angular-mocks.js",
		"bower_components/angular-mocks/angular-mocks.js",
	)
}

// bower_components_angular_mocks_bower_json reads file data from disk. It returns an error on failure.
func bower_components_angular_mocks_bower_json() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-mocks/bower.json",
		"bower_components/angular-mocks/bower.json",
	)
}

// bower_components_angular_mocks_readme_md reads file data from disk. It returns an error on failure.
func bower_components_angular_mocks_readme_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-mocks/README.md",
		"bower_components/angular-mocks/README.md",
	)
}

// bower_components_angular_resource_angular_resource_js reads file data from disk. It returns an error on failure.
func bower_components_angular_resource_angular_resource_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-resource/angular-resource.js",
		"bower_components/angular-resource/angular-resource.js",
	)
}

// bower_components_angular_resource_angular_resource_min_js reads file data from disk. It returns an error on failure.
func bower_components_angular_resource_angular_resource_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-resource/angular-resource.min.js",
		"bower_components/angular-resource/angular-resource.min.js",
	)
}

// bower_components_angular_resource_angular_resource_min_js_map reads file data from disk. It returns an error on failure.
func bower_components_angular_resource_angular_resource_min_js_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-resource/angular-resource.min.js.map",
		"bower_components/angular-resource/angular-resource.min.js.map",
	)
}

// bower_components_angular_resource_bower_json reads file data from disk. It returns an error on failure.
func bower_components_angular_resource_bower_json() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-resource/bower.json",
		"bower_components/angular-resource/bower.json",
	)
}

// bower_components_angular_resource_readme_md reads file data from disk. It returns an error on failure.
func bower_components_angular_resource_readme_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-resource/README.md",
		"bower_components/angular-resource/README.md",
	)
}

// bower_components_angular_route_angular_route_js reads file data from disk. It returns an error on failure.
func bower_components_angular_route_angular_route_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-route/angular-route.js",
		"bower_components/angular-route/angular-route.js",
	)
}

// bower_components_angular_route_angular_route_min_js reads file data from disk. It returns an error on failure.
func bower_components_angular_route_angular_route_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-route/angular-route.min.js",
		"bower_components/angular-route/angular-route.min.js",
	)
}

// bower_components_angular_route_angular_route_min_js_map reads file data from disk. It returns an error on failure.
func bower_components_angular_route_angular_route_min_js_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-route/angular-route.min.js.map",
		"bower_components/angular-route/angular-route.min.js.map",
	)
}

// bower_components_angular_route_bower_json reads file data from disk. It returns an error on failure.
func bower_components_angular_route_bower_json() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-route/bower.json",
		"bower_components/angular-route/bower.json",
	)
}

// bower_components_angular_route_readme_md reads file data from disk. It returns an error on failure.
func bower_components_angular_route_readme_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/angular-route/README.md",
		"bower_components/angular-route/README.md",
	)
}

// bower_components_html5_boilerplate_404_html reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_404_html() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/404.html",
		"bower_components/html5-boilerplate/404.html",
	)
}

// bower_components_html5_boilerplate_apple_touch_icon_precomposed_png reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_apple_touch_icon_precomposed_png() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/apple-touch-icon-precomposed.png",
		"bower_components/html5-boilerplate/apple-touch-icon-precomposed.png",
	)
}

// bower_components_html5_boilerplate_changelog_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_changelog_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/CHANGELOG.md",
		"bower_components/html5-boilerplate/CHANGELOG.md",
	)
}

// bower_components_html5_boilerplate_contributing_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_contributing_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/CONTRIBUTING.md",
		"bower_components/html5-boilerplate/CONTRIBUTING.md",
	)
}

// bower_components_html5_boilerplate_crossdomain_xml reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_crossdomain_xml() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/crossdomain.xml",
		"bower_components/html5-boilerplate/crossdomain.xml",
	)
}

// bower_components_html5_boilerplate_css_main_css reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_css_main_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/css/main.css",
		"bower_components/html5-boilerplate/css/main.css",
	)
}

// bower_components_html5_boilerplate_css_normalize_css reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_css_normalize_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/css/normalize.css",
		"bower_components/html5-boilerplate/css/normalize.css",
	)
}

// bower_components_html5_boilerplate_doc_crossdomain_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_crossdomain_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/crossdomain.md",
		"bower_components/html5-boilerplate/doc/crossdomain.md",
	)
}

// bower_components_html5_boilerplate_doc_css_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_css_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/css.md",
		"bower_components/html5-boilerplate/doc/css.md",
	)
}

// bower_components_html5_boilerplate_doc_extend_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_extend_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/extend.md",
		"bower_components/html5-boilerplate/doc/extend.md",
	)
}

// bower_components_html5_boilerplate_doc_faq_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_faq_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/faq.md",
		"bower_components/html5-boilerplate/doc/faq.md",
	)
}

// bower_components_html5_boilerplate_doc_html_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_html_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/html.md",
		"bower_components/html5-boilerplate/doc/html.md",
	)
}

// bower_components_html5_boilerplate_doc_js_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_js_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/js.md",
		"bower_components/html5-boilerplate/doc/js.md",
	)
}

// bower_components_html5_boilerplate_doc_misc_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_misc_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/misc.md",
		"bower_components/html5-boilerplate/doc/misc.md",
	)
}

// bower_components_html5_boilerplate_doc_toc_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_toc_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/TOC.md",
		"bower_components/html5-boilerplate/doc/TOC.md",
	)
}

// bower_components_html5_boilerplate_doc_usage_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_doc_usage_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/doc/usage.md",
		"bower_components/html5-boilerplate/doc/usage.md",
	)
}

// bower_components_html5_boilerplate_favicon_ico reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_favicon_ico() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/favicon.ico",
		"bower_components/html5-boilerplate/favicon.ico",
	)
}

// bower_components_html5_boilerplate_humans_txt reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_humans_txt() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/humans.txt",
		"bower_components/html5-boilerplate/humans.txt",
	)
}

// bower_components_html5_boilerplate_index_html reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_index_html() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/index.html",
		"bower_components/html5-boilerplate/index.html",
	)
}

// bower_components_html5_boilerplate_js_main_js reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_js_main_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/js/main.js",
		"bower_components/html5-boilerplate/js/main.js",
	)
}

// bower_components_html5_boilerplate_js_plugins_js reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_js_plugins_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/js/plugins.js",
		"bower_components/html5-boilerplate/js/plugins.js",
	)
}

// bower_components_html5_boilerplate_js_vendor_jquery_1_10_2_min_js reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_js_vendor_jquery_1_10_2_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/js/vendor/jquery-1.10.2.min.js",
		"bower_components/html5-boilerplate/js/vendor/jquery-1.10.2.min.js",
	)
}

// bower_components_html5_boilerplate_js_vendor_modernizr_2_6_2_min_js reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_js_vendor_modernizr_2_6_2_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/js/vendor/modernizr-2.6.2.min.js",
		"bower_components/html5-boilerplate/js/vendor/modernizr-2.6.2.min.js",
	)
}

// bower_components_html5_boilerplate_license_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_license_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/LICENSE.md",
		"bower_components/html5-boilerplate/LICENSE.md",
	)
}

// bower_components_html5_boilerplate_readme_md reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_readme_md() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/README.md",
		"bower_components/html5-boilerplate/README.md",
	)
}

// bower_components_html5_boilerplate_robots_txt reads file data from disk. It returns an error on failure.
func bower_components_html5_boilerplate_robots_txt() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/html5-boilerplate/robots.txt",
		"bower_components/html5-boilerplate/robots.txt",
	)
}

// css_app_css reads file data from disk. It returns an error on failure.
func css_app_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/app.css",
		"css/app.css",
	)
}

// index_async_html reads file data from disk. It returns an error on failure.
func index_async_html() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/index-async.html",
		"index-async.html",
	)
}

// index_html reads file data from disk. It returns an error on failure.
func index_html() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/index.html",
		"index.html",
	)
}

// js_app_js reads file data from disk. It returns an error on failure.
func js_app_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/app.js",
		"js/app.js",
	)
}

// js_controllers_js reads file data from disk. It returns an error on failure.
func js_controllers_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/controllers.js",
		"js/controllers.js",
	)
}

// js_directives_js reads file data from disk. It returns an error on failure.
func js_directives_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/directives.js",
		"js/directives.js",
	)
}

// js_filters_js reads file data from disk. It returns an error on failure.
func js_filters_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/filters.js",
		"js/filters.js",
	)
}

// js_services_js reads file data from disk. It returns an error on failure.
func js_services_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/services.js",
		"js/services.js",
	)
}

// npm_debug_log reads file data from disk. It returns an error on failure.
func npm_debug_log() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/npm-debug.log",
		"npm-debug.log",
	)
}

// partials_partial1_html reads file data from disk. It returns an error on failure.
func partials_partial1_html() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/partials/partial1.html",
		"partials/partial1.html",
	)
}

// partials_partial2_html reads file data from disk. It returns an error on failure.
func partials_partial2_html() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/partials/partial2.html",
		"partials/partial2.html",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"bower_components/angular/angular-csp.css":                            bower_components_angular_angular_csp_css,
	"bower_components/angular/angular.js":                                 bower_components_angular_angular_js,
	"bower_components/angular/angular.min.js":                             bower_components_angular_angular_min_js,
	"bower_components/angular/angular.min.js.gzip":                        bower_components_angular_angular_min_js_gzip,
	"bower_components/angular/angular.min.js.map":                         bower_components_angular_angular_min_js_map,
	"bower_components/angular/bower.json":                                 bower_components_angular_bower_json,
	"bower_components/angular/README.md":                                  bower_components_angular_readme_md,
	"bower_components/angular-loader/angular-loader.js":                   bower_components_angular_loader_angular_loader_js,
	"bower_components/angular-loader/angular-loader.min.js":               bower_components_angular_loader_angular_loader_min_js,
	"bower_components/angular-loader/angular-loader.min.js.map":           bower_components_angular_loader_angular_loader_min_js_map,
	"bower_components/angular-loader/bower.json":                          bower_components_angular_loader_bower_json,
	"bower_components/angular-loader/README.md":                           bower_components_angular_loader_readme_md,
	"bower_components/angular-mocks/angular-mocks.js":                     bower_components_angular_mocks_angular_mocks_js,
	"bower_components/angular-mocks/bower.json":                           bower_components_angular_mocks_bower_json,
	"bower_components/angular-mocks/README.md":                            bower_components_angular_mocks_readme_md,
	"bower_components/angular-resource/angular-resource.js":               bower_components_angular_resource_angular_resource_js,
	"bower_components/angular-resource/angular-resource.min.js":           bower_components_angular_resource_angular_resource_min_js,
	"bower_components/angular-resource/angular-resource.min.js.map":       bower_components_angular_resource_angular_resource_min_js_map,
	"bower_components/angular-resource/bower.json":                        bower_components_angular_resource_bower_json,
	"bower_components/angular-resource/README.md":                         bower_components_angular_resource_readme_md,
	"bower_components/angular-route/angular-route.js":                     bower_components_angular_route_angular_route_js,
	"bower_components/angular-route/angular-route.min.js":                 bower_components_angular_route_angular_route_min_js,
	"bower_components/angular-route/angular-route.min.js.map":             bower_components_angular_route_angular_route_min_js_map,
	"bower_components/angular-route/bower.json":                           bower_components_angular_route_bower_json,
	"bower_components/angular-route/README.md":                            bower_components_angular_route_readme_md,
	"bower_components/html5-boilerplate/404.html":                         bower_components_html5_boilerplate_404_html,
	"bower_components/html5-boilerplate/apple-touch-icon-precomposed.png": bower_components_html5_boilerplate_apple_touch_icon_precomposed_png,
	"bower_components/html5-boilerplate/CHANGELOG.md":                     bower_components_html5_boilerplate_changelog_md,
	"bower_components/html5-boilerplate/CONTRIBUTING.md":                  bower_components_html5_boilerplate_contributing_md,
	"bower_components/html5-boilerplate/crossdomain.xml":                  bower_components_html5_boilerplate_crossdomain_xml,
	"bower_components/html5-boilerplate/css/main.css":                     bower_components_html5_boilerplate_css_main_css,
	"bower_components/html5-boilerplate/css/normalize.css":                bower_components_html5_boilerplate_css_normalize_css,
	"bower_components/html5-boilerplate/doc/crossdomain.md":               bower_components_html5_boilerplate_doc_crossdomain_md,
	"bower_components/html5-boilerplate/doc/css.md":                       bower_components_html5_boilerplate_doc_css_md,
	"bower_components/html5-boilerplate/doc/extend.md":                    bower_components_html5_boilerplate_doc_extend_md,
	"bower_components/html5-boilerplate/doc/faq.md":                       bower_components_html5_boilerplate_doc_faq_md,
	"bower_components/html5-boilerplate/doc/html.md":                      bower_components_html5_boilerplate_doc_html_md,
	"bower_components/html5-boilerplate/doc/js.md":                        bower_components_html5_boilerplate_doc_js_md,
	"bower_components/html5-boilerplate/doc/misc.md":                      bower_components_html5_boilerplate_doc_misc_md,
	"bower_components/html5-boilerplate/doc/TOC.md":                       bower_components_html5_boilerplate_doc_toc_md,
	"bower_components/html5-boilerplate/doc/usage.md":                     bower_components_html5_boilerplate_doc_usage_md,
	"bower_components/html5-boilerplate/favicon.ico":                      bower_components_html5_boilerplate_favicon_ico,
	"bower_components/html5-boilerplate/humans.txt":                       bower_components_html5_boilerplate_humans_txt,
	"bower_components/html5-boilerplate/index.html":                       bower_components_html5_boilerplate_index_html,
	"bower_components/html5-boilerplate/js/main.js":                       bower_components_html5_boilerplate_js_main_js,
	"bower_components/html5-boilerplate/js/plugins.js":                    bower_components_html5_boilerplate_js_plugins_js,
	"bower_components/html5-boilerplate/js/vendor/jquery-1.10.2.min.js":   bower_components_html5_boilerplate_js_vendor_jquery_1_10_2_min_js,
	"bower_components/html5-boilerplate/js/vendor/modernizr-2.6.2.min.js": bower_components_html5_boilerplate_js_vendor_modernizr_2_6_2_min_js,
	"bower_components/html5-boilerplate/LICENSE.md":                       bower_components_html5_boilerplate_license_md,
	"bower_components/html5-boilerplate/README.md":                        bower_components_html5_boilerplate_readme_md,
	"bower_components/html5-boilerplate/robots.txt":                       bower_components_html5_boilerplate_robots_txt,
	"css/app.css":                                                         css_app_css,
	"index-async.html":                                                    index_async_html,
	"index.html":                                                          index_html,
	"js/app.js":                                                           js_app_js,
	"js/controllers.js":                                                   js_controllers_js,
	"js/directives.js":                                                    js_directives_js,
	"js/filters.js":                                                       js_filters_js,
	"js/services.js":                                                      js_services_js,
	"npm-debug.log":                                                       npm_debug_log,
	"partials/partial1.html":                                              partials_partial1_html,
	"partials/partial2.html":                                              partials_partial2_html,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"bower_components": &_bintree_t{nil, map[string]*_bintree_t{
		"angular": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-csp.css":     &_bintree_t{bower_components_angular_angular_csp_css, map[string]*_bintree_t{}},
			"angular.js":          &_bintree_t{bower_components_angular_angular_js, map[string]*_bintree_t{}},
			"angular.min.js":      &_bintree_t{bower_components_angular_angular_min_js, map[string]*_bintree_t{}},
			"angular.min.js.gzip": &_bintree_t{bower_components_angular_angular_min_js_gzip, map[string]*_bintree_t{}},
			"angular.min.js.map":  &_bintree_t{bower_components_angular_angular_min_js_map, map[string]*_bintree_t{}},
			"bower.json":          &_bintree_t{bower_components_angular_bower_json, map[string]*_bintree_t{}},
			"README.md":           &_bintree_t{bower_components_angular_readme_md, map[string]*_bintree_t{}},
		}},
		"angular-loader": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-loader.min.js":     &_bintree_t{bower_components_angular_loader_angular_loader_min_js, map[string]*_bintree_t{}},
			"angular-loader.min.js.map": &_bintree_t{bower_components_angular_loader_angular_loader_min_js_map, map[string]*_bintree_t{}},
			"bower.json":                &_bintree_t{bower_components_angular_loader_bower_json, map[string]*_bintree_t{}},
			"README.md":                 &_bintree_t{bower_components_angular_loader_readme_md, map[string]*_bintree_t{}},
			"angular-loader.js":         &_bintree_t{bower_components_angular_loader_angular_loader_js, map[string]*_bintree_t{}},
		}},
		"angular-mocks": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-mocks.js": &_bintree_t{bower_components_angular_mocks_angular_mocks_js, map[string]*_bintree_t{}},
			"bower.json":       &_bintree_t{bower_components_angular_mocks_bower_json, map[string]*_bintree_t{}},
			"README.md":        &_bintree_t{bower_components_angular_mocks_readme_md, map[string]*_bintree_t{}},
		}},
		"angular-resource": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-resource.js":         &_bintree_t{bower_components_angular_resource_angular_resource_js, map[string]*_bintree_t{}},
			"angular-resource.min.js":     &_bintree_t{bower_components_angular_resource_angular_resource_min_js, map[string]*_bintree_t{}},
			"angular-resource.min.js.map": &_bintree_t{bower_components_angular_resource_angular_resource_min_js_map, map[string]*_bintree_t{}},
			"bower.json":                  &_bintree_t{bower_components_angular_resource_bower_json, map[string]*_bintree_t{}},
			"README.md":                   &_bintree_t{bower_components_angular_resource_readme_md, map[string]*_bintree_t{}},
		}},
		"angular-route": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-route.min.js.map": &_bintree_t{bower_components_angular_route_angular_route_min_js_map, map[string]*_bintree_t{}},
			"bower.json":               &_bintree_t{bower_components_angular_route_bower_json, map[string]*_bintree_t{}},
			"README.md":                &_bintree_t{bower_components_angular_route_readme_md, map[string]*_bintree_t{}},
			"angular-route.js":         &_bintree_t{bower_components_angular_route_angular_route_js, map[string]*_bintree_t{}},
			"angular-route.min.js":     &_bintree_t{bower_components_angular_route_angular_route_min_js, map[string]*_bintree_t{}},
		}},
		"html5-boilerplate": &_bintree_t{nil, map[string]*_bintree_t{
			"crossdomain.xml": &_bintree_t{bower_components_html5_boilerplate_crossdomain_xml, map[string]*_bintree_t{}},
			"css": &_bintree_t{nil, map[string]*_bintree_t{
				"main.css":      &_bintree_t{bower_components_html5_boilerplate_css_main_css, map[string]*_bintree_t{}},
				"normalize.css": &_bintree_t{bower_components_html5_boilerplate_css_normalize_css, map[string]*_bintree_t{}},
			}},
			"humans.txt":      &_bintree_t{bower_components_html5_boilerplate_humans_txt, map[string]*_bintree_t{}},
			"README.md":       &_bintree_t{bower_components_html5_boilerplate_readme_md, map[string]*_bintree_t{}},
			"404.html":        &_bintree_t{bower_components_html5_boilerplate_404_html, map[string]*_bintree_t{}},
			"CONTRIBUTING.md": &_bintree_t{bower_components_html5_boilerplate_contributing_md, map[string]*_bintree_t{}},
			"doc": &_bintree_t{nil, map[string]*_bintree_t{
				"css.md":         &_bintree_t{bower_components_html5_boilerplate_doc_css_md, map[string]*_bintree_t{}},
				"extend.md":      &_bintree_t{bower_components_html5_boilerplate_doc_extend_md, map[string]*_bintree_t{}},
				"faq.md":         &_bintree_t{bower_components_html5_boilerplate_doc_faq_md, map[string]*_bintree_t{}},
				"html.md":        &_bintree_t{bower_components_html5_boilerplate_doc_html_md, map[string]*_bintree_t{}},
				"js.md":          &_bintree_t{bower_components_html5_boilerplate_doc_js_md, map[string]*_bintree_t{}},
				"usage.md":       &_bintree_t{bower_components_html5_boilerplate_doc_usage_md, map[string]*_bintree_t{}},
				"crossdomain.md": &_bintree_t{bower_components_html5_boilerplate_doc_crossdomain_md, map[string]*_bintree_t{}},
				"misc.md":        &_bintree_t{bower_components_html5_boilerplate_doc_misc_md, map[string]*_bintree_t{}},
				"TOC.md":         &_bintree_t{bower_components_html5_boilerplate_doc_toc_md, map[string]*_bintree_t{}},
			}},
			"js": &_bintree_t{nil, map[string]*_bintree_t{
				"main.js":    &_bintree_t{bower_components_html5_boilerplate_js_main_js, map[string]*_bintree_t{}},
				"plugins.js": &_bintree_t{bower_components_html5_boilerplate_js_plugins_js, map[string]*_bintree_t{}},
				"vendor": &_bintree_t{nil, map[string]*_bintree_t{
					"jquery-1.10.2.min.js":   &_bintree_t{bower_components_html5_boilerplate_js_vendor_jquery_1_10_2_min_js, map[string]*_bintree_t{}},
					"modernizr-2.6.2.min.js": &_bintree_t{bower_components_html5_boilerplate_js_vendor_modernizr_2_6_2_min_js, map[string]*_bintree_t{}},
				}},
			}},
			"CHANGELOG.md":                     &_bintree_t{bower_components_html5_boilerplate_changelog_md, map[string]*_bintree_t{}},
			"LICENSE.md":                       &_bintree_t{bower_components_html5_boilerplate_license_md, map[string]*_bintree_t{}},
			"robots.txt":                       &_bintree_t{bower_components_html5_boilerplate_robots_txt, map[string]*_bintree_t{}},
			"apple-touch-icon-precomposed.png": &_bintree_t{bower_components_html5_boilerplate_apple_touch_icon_precomposed_png, map[string]*_bintree_t{}},
			"favicon.ico":                      &_bintree_t{bower_components_html5_boilerplate_favicon_ico, map[string]*_bintree_t{}},
			"index.html":                       &_bintree_t{bower_components_html5_boilerplate_index_html, map[string]*_bintree_t{}},
		}},
	}},
	"css": &_bintree_t{nil, map[string]*_bintree_t{
		"app.css": &_bintree_t{css_app_css, map[string]*_bintree_t{}},
	}},
	"index-async.html": &_bintree_t{index_async_html, map[string]*_bintree_t{}},
	"index.html":       &_bintree_t{index_html, map[string]*_bintree_t{}},
	"js": &_bintree_t{nil, map[string]*_bintree_t{
		"directives.js":  &_bintree_t{js_directives_js, map[string]*_bintree_t{}},
		"filters.js":     &_bintree_t{js_filters_js, map[string]*_bintree_t{}},
		"services.js":    &_bintree_t{js_services_js, map[string]*_bintree_t{}},
		"app.js":         &_bintree_t{js_app_js, map[string]*_bintree_t{}},
		"controllers.js": &_bintree_t{js_controllers_js, map[string]*_bintree_t{}},
	}},
	"npm-debug.log": &_bintree_t{npm_debug_log, map[string]*_bintree_t{}},
	"partials": &_bintree_t{nil, map[string]*_bintree_t{
		"partial1.html": &_bintree_t{partials_partial1_html, map[string]*_bintree_t{}},
		"partial2.html": &_bintree_t{partials_partial2_html, map[string]*_bintree_t{}},
	}},
}}
