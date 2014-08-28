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

// bower_components_jquery_bower_json reads file data from disk. It returns an error on failure.
func bower_components_jquery_bower_json() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/bower.json",
		"bower_components/jquery/bower.json",
	)
}

// bower_components_jquery_dist_jquery_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_dist_jquery_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/dist/jquery.js",
		"bower_components/jquery/dist/jquery.js",
	)
}

// bower_components_jquery_dist_jquery_min_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_dist_jquery_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/dist/jquery.min.js",
		"bower_components/jquery/dist/jquery.min.js",
	)
}

// bower_components_jquery_dist_jquery_min_map reads file data from disk. It returns an error on failure.
func bower_components_jquery_dist_jquery_min_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/dist/jquery.min.map",
		"bower_components/jquery/dist/jquery.min.map",
	)
}

// bower_components_jquery_mit_license_txt reads file data from disk. It returns an error on failure.
func bower_components_jquery_mit_license_txt() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/MIT-LICENSE.txt",
		"bower_components/jquery/MIT-LICENSE.txt",
	)
}

// bower_components_jquery_src_ajax_jsonp_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_jsonp_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/jsonp.js",
		"bower_components/jquery/src/ajax/jsonp.js",
	)
}

// bower_components_jquery_src_ajax_load_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_load_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/load.js",
		"bower_components/jquery/src/ajax/load.js",
	)
}

// bower_components_jquery_src_ajax_parsejson_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_parsejson_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/parseJSON.js",
		"bower_components/jquery/src/ajax/parseJSON.js",
	)
}

// bower_components_jquery_src_ajax_parsexml_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_parsexml_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/parseXML.js",
		"bower_components/jquery/src/ajax/parseXML.js",
	)
}

// bower_components_jquery_src_ajax_script_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_script_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/script.js",
		"bower_components/jquery/src/ajax/script.js",
	)
}

// bower_components_jquery_src_ajax_var_nonce_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_var_nonce_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/var/nonce.js",
		"bower_components/jquery/src/ajax/var/nonce.js",
	)
}

// bower_components_jquery_src_ajax_var_rquery_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_var_rquery_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/var/rquery.js",
		"bower_components/jquery/src/ajax/var/rquery.js",
	)
}

// bower_components_jquery_src_ajax_xhr_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_xhr_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax/xhr.js",
		"bower_components/jquery/src/ajax/xhr.js",
	)
}

// bower_components_jquery_src_ajax_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_ajax_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/ajax.js",
		"bower_components/jquery/src/ajax.js",
	)
}

// bower_components_jquery_src_attributes_attr_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_attributes_attr_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/attributes/attr.js",
		"bower_components/jquery/src/attributes/attr.js",
	)
}

// bower_components_jquery_src_attributes_classes_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_attributes_classes_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/attributes/classes.js",
		"bower_components/jquery/src/attributes/classes.js",
	)
}

// bower_components_jquery_src_attributes_prop_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_attributes_prop_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/attributes/prop.js",
		"bower_components/jquery/src/attributes/prop.js",
	)
}

// bower_components_jquery_src_attributes_support_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_attributes_support_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/attributes/support.js",
		"bower_components/jquery/src/attributes/support.js",
	)
}

// bower_components_jquery_src_attributes_val_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_attributes_val_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/attributes/val.js",
		"bower_components/jquery/src/attributes/val.js",
	)
}

// bower_components_jquery_src_attributes_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_attributes_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/attributes.js",
		"bower_components/jquery/src/attributes.js",
	)
}

// bower_components_jquery_src_callbacks_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_callbacks_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/callbacks.js",
		"bower_components/jquery/src/callbacks.js",
	)
}

// bower_components_jquery_src_core_access_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_core_access_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/core/access.js",
		"bower_components/jquery/src/core/access.js",
	)
}

// bower_components_jquery_src_core_init_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_core_init_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/core/init.js",
		"bower_components/jquery/src/core/init.js",
	)
}

// bower_components_jquery_src_core_parsehtml_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_core_parsehtml_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/core/parseHTML.js",
		"bower_components/jquery/src/core/parseHTML.js",
	)
}

// bower_components_jquery_src_core_ready_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_core_ready_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/core/ready.js",
		"bower_components/jquery/src/core/ready.js",
	)
}

// bower_components_jquery_src_core_var_rsingletag_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_core_var_rsingletag_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/core/var/rsingleTag.js",
		"bower_components/jquery/src/core/var/rsingleTag.js",
	)
}

// bower_components_jquery_src_core_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_core_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/core.js",
		"bower_components/jquery/src/core.js",
	)
}

// bower_components_jquery_src_css_addgethookif_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_addgethookif_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/addGetHookIf.js",
		"bower_components/jquery/src/css/addGetHookIf.js",
	)
}

// bower_components_jquery_src_css_curcss_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_curcss_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/curCSS.js",
		"bower_components/jquery/src/css/curCSS.js",
	)
}

// bower_components_jquery_src_css_defaultdisplay_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_defaultdisplay_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/defaultDisplay.js",
		"bower_components/jquery/src/css/defaultDisplay.js",
	)
}

// bower_components_jquery_src_css_hiddenvisibleselectors_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_hiddenvisibleselectors_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/hiddenVisibleSelectors.js",
		"bower_components/jquery/src/css/hiddenVisibleSelectors.js",
	)
}

// bower_components_jquery_src_css_support_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_support_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/support.js",
		"bower_components/jquery/src/css/support.js",
	)
}

// bower_components_jquery_src_css_swap_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_swap_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/swap.js",
		"bower_components/jquery/src/css/swap.js",
	)
}

// bower_components_jquery_src_css_var_cssexpand_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_var_cssexpand_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/var/cssExpand.js",
		"bower_components/jquery/src/css/var/cssExpand.js",
	)
}

// bower_components_jquery_src_css_var_getstyles_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_var_getstyles_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/var/getStyles.js",
		"bower_components/jquery/src/css/var/getStyles.js",
	)
}

// bower_components_jquery_src_css_var_ishidden_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_var_ishidden_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/var/isHidden.js",
		"bower_components/jquery/src/css/var/isHidden.js",
	)
}

// bower_components_jquery_src_css_var_rmargin_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_var_rmargin_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/var/rmargin.js",
		"bower_components/jquery/src/css/var/rmargin.js",
	)
}

// bower_components_jquery_src_css_var_rnumnonpx_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_var_rnumnonpx_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css/var/rnumnonpx.js",
		"bower_components/jquery/src/css/var/rnumnonpx.js",
	)
}

// bower_components_jquery_src_css_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_css_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/css.js",
		"bower_components/jquery/src/css.js",
	)
}

// bower_components_jquery_src_data_accepts_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_data_accepts_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/data/accepts.js",
		"bower_components/jquery/src/data/accepts.js",
	)
}

// bower_components_jquery_src_data_data_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_data_data_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/data/Data.js",
		"bower_components/jquery/src/data/Data.js",
	)
}

// bower_components_jquery_src_data_var_data_priv_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_data_var_data_priv_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/data/var/data_priv.js",
		"bower_components/jquery/src/data/var/data_priv.js",
	)
}

// bower_components_jquery_src_data_var_data_user_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_data_var_data_user_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/data/var/data_user.js",
		"bower_components/jquery/src/data/var/data_user.js",
	)
}

// bower_components_jquery_src_data_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_data_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/data.js",
		"bower_components/jquery/src/data.js",
	)
}

// bower_components_jquery_src_deferred_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_deferred_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/deferred.js",
		"bower_components/jquery/src/deferred.js",
	)
}

// bower_components_jquery_src_deprecated_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_deprecated_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/deprecated.js",
		"bower_components/jquery/src/deprecated.js",
	)
}

// bower_components_jquery_src_dimensions_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_dimensions_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/dimensions.js",
		"bower_components/jquery/src/dimensions.js",
	)
}

// bower_components_jquery_src_effects_animatedselector_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_effects_animatedselector_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/effects/animatedSelector.js",
		"bower_components/jquery/src/effects/animatedSelector.js",
	)
}

// bower_components_jquery_src_effects_tween_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_effects_tween_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/effects/Tween.js",
		"bower_components/jquery/src/effects/Tween.js",
	)
}

// bower_components_jquery_src_effects_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_effects_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/effects.js",
		"bower_components/jquery/src/effects.js",
	)
}

// bower_components_jquery_src_event_alias_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_event_alias_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/event/alias.js",
		"bower_components/jquery/src/event/alias.js",
	)
}

// bower_components_jquery_src_event_support_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_event_support_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/event/support.js",
		"bower_components/jquery/src/event/support.js",
	)
}

// bower_components_jquery_src_event_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_event_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/event.js",
		"bower_components/jquery/src/event.js",
	)
}

// bower_components_jquery_src_exports_amd_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_exports_amd_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/exports/amd.js",
		"bower_components/jquery/src/exports/amd.js",
	)
}

// bower_components_jquery_src_exports_global_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_exports_global_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/exports/global.js",
		"bower_components/jquery/src/exports/global.js",
	)
}

// bower_components_jquery_src_intro_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_intro_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/intro.js",
		"bower_components/jquery/src/intro.js",
	)
}

// bower_components_jquery_src_jquery_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_jquery_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/jquery.js",
		"bower_components/jquery/src/jquery.js",
	)
}

// bower_components_jquery_src_manipulation_evalurl_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_manipulation_evalurl_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/manipulation/_evalUrl.js",
		"bower_components/jquery/src/manipulation/_evalUrl.js",
	)
}

// bower_components_jquery_src_manipulation_support_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_manipulation_support_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/manipulation/support.js",
		"bower_components/jquery/src/manipulation/support.js",
	)
}

// bower_components_jquery_src_manipulation_var_rcheckabletype_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_manipulation_var_rcheckabletype_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/manipulation/var/rcheckableType.js",
		"bower_components/jquery/src/manipulation/var/rcheckableType.js",
	)
}

// bower_components_jquery_src_manipulation_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_manipulation_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/manipulation.js",
		"bower_components/jquery/src/manipulation.js",
	)
}

// bower_components_jquery_src_offset_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_offset_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/offset.js",
		"bower_components/jquery/src/offset.js",
	)
}

// bower_components_jquery_src_outro_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_outro_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/outro.js",
		"bower_components/jquery/src/outro.js",
	)
}

// bower_components_jquery_src_queue_delay_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_queue_delay_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/queue/delay.js",
		"bower_components/jquery/src/queue/delay.js",
	)
}

// bower_components_jquery_src_queue_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_queue_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/queue.js",
		"bower_components/jquery/src/queue.js",
	)
}

// bower_components_jquery_src_selector_native_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_selector_native_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/selector-native.js",
		"bower_components/jquery/src/selector-native.js",
	)
}

// bower_components_jquery_src_selector_sizzle_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_selector_sizzle_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/selector-sizzle.js",
		"bower_components/jquery/src/selector-sizzle.js",
	)
}

// bower_components_jquery_src_selector_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_selector_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/selector.js",
		"bower_components/jquery/src/selector.js",
	)
}

// bower_components_jquery_src_serialize_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_serialize_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/serialize.js",
		"bower_components/jquery/src/serialize.js",
	)
}

// bower_components_jquery_src_sizzle_dist_sizzle_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_sizzle_dist_sizzle_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/sizzle/dist/sizzle.js",
		"bower_components/jquery/src/sizzle/dist/sizzle.js",
	)
}

// bower_components_jquery_src_sizzle_dist_sizzle_min_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_sizzle_dist_sizzle_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/sizzle/dist/sizzle.min.js",
		"bower_components/jquery/src/sizzle/dist/sizzle.min.js",
	)
}

// bower_components_jquery_src_sizzle_dist_sizzle_min_map reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_sizzle_dist_sizzle_min_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/sizzle/dist/sizzle.min.map",
		"bower_components/jquery/src/sizzle/dist/sizzle.min.map",
	)
}

// bower_components_jquery_src_traversing_findfilter_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_traversing_findfilter_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/traversing/findFilter.js",
		"bower_components/jquery/src/traversing/findFilter.js",
	)
}

// bower_components_jquery_src_traversing_var_rneedscontext_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_traversing_var_rneedscontext_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/traversing/var/rneedsContext.js",
		"bower_components/jquery/src/traversing/var/rneedsContext.js",
	)
}

// bower_components_jquery_src_traversing_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_traversing_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/traversing.js",
		"bower_components/jquery/src/traversing.js",
	)
}

// bower_components_jquery_src_var_arr_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_arr_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/arr.js",
		"bower_components/jquery/src/var/arr.js",
	)
}

// bower_components_jquery_src_var_class2type_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_class2type_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/class2type.js",
		"bower_components/jquery/src/var/class2type.js",
	)
}

// bower_components_jquery_src_var_concat_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_concat_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/concat.js",
		"bower_components/jquery/src/var/concat.js",
	)
}

// bower_components_jquery_src_var_hasown_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_hasown_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/hasOwn.js",
		"bower_components/jquery/src/var/hasOwn.js",
	)
}

// bower_components_jquery_src_var_indexof_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_indexof_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/indexOf.js",
		"bower_components/jquery/src/var/indexOf.js",
	)
}

// bower_components_jquery_src_var_pnum_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_pnum_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/pnum.js",
		"bower_components/jquery/src/var/pnum.js",
	)
}

// bower_components_jquery_src_var_push_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_push_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/push.js",
		"bower_components/jquery/src/var/push.js",
	)
}

// bower_components_jquery_src_var_rnotwhite_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_rnotwhite_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/rnotwhite.js",
		"bower_components/jquery/src/var/rnotwhite.js",
	)
}

// bower_components_jquery_src_var_slice_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_slice_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/slice.js",
		"bower_components/jquery/src/var/slice.js",
	)
}

// bower_components_jquery_src_var_strundefined_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_strundefined_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/strundefined.js",
		"bower_components/jquery/src/var/strundefined.js",
	)
}

// bower_components_jquery_src_var_support_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_support_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/support.js",
		"bower_components/jquery/src/var/support.js",
	)
}

// bower_components_jquery_src_var_tostring_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_var_tostring_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/var/toString.js",
		"bower_components/jquery/src/var/toString.js",
	)
}

// bower_components_jquery_src_wrap_js reads file data from disk. It returns an error on failure.
func bower_components_jquery_src_wrap_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/bower_components/jquery/src/wrap.js",
		"bower_components/jquery/src/wrap.js",
	)
}

// css_app_css reads file data from disk. It returns an error on failure.
func css_app_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/app.css",
		"css/app.css",
	)
}

// css_bootstrap_css_bootstrap_theme_css reads file data from disk. It returns an error on failure.
func css_bootstrap_css_bootstrap_theme_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/css/bootstrap-theme.css",
		"css/bootstrap/css/bootstrap-theme.css",
	)
}

// css_bootstrap_css_bootstrap_theme_css_map reads file data from disk. It returns an error on failure.
func css_bootstrap_css_bootstrap_theme_css_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/css/bootstrap-theme.css.map",
		"css/bootstrap/css/bootstrap-theme.css.map",
	)
}

// css_bootstrap_css_bootstrap_theme_min_css reads file data from disk. It returns an error on failure.
func css_bootstrap_css_bootstrap_theme_min_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/css/bootstrap-theme.min.css",
		"css/bootstrap/css/bootstrap-theme.min.css",
	)
}

// css_bootstrap_css_bootstrap_css reads file data from disk. It returns an error on failure.
func css_bootstrap_css_bootstrap_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/css/bootstrap.css",
		"css/bootstrap/css/bootstrap.css",
	)
}

// css_bootstrap_css_bootstrap_css_map reads file data from disk. It returns an error on failure.
func css_bootstrap_css_bootstrap_css_map() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/css/bootstrap.css.map",
		"css/bootstrap/css/bootstrap.css.map",
	)
}

// css_bootstrap_css_bootstrap_min_css reads file data from disk. It returns an error on failure.
func css_bootstrap_css_bootstrap_min_css() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/css/bootstrap.min.css",
		"css/bootstrap/css/bootstrap.min.css",
	)
}

// css_bootstrap_fonts_glyphicons_halflings_regular_eot reads file data from disk. It returns an error on failure.
func css_bootstrap_fonts_glyphicons_halflings_regular_eot() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/fonts/glyphicons-halflings-regular.eot",
		"css/bootstrap/fonts/glyphicons-halflings-regular.eot",
	)
}

// css_bootstrap_fonts_glyphicons_halflings_regular_svg reads file data from disk. It returns an error on failure.
func css_bootstrap_fonts_glyphicons_halflings_regular_svg() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/fonts/glyphicons-halflings-regular.svg",
		"css/bootstrap/fonts/glyphicons-halflings-regular.svg",
	)
}

// css_bootstrap_fonts_glyphicons_halflings_regular_ttf reads file data from disk. It returns an error on failure.
func css_bootstrap_fonts_glyphicons_halflings_regular_ttf() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/fonts/glyphicons-halflings-regular.ttf",
		"css/bootstrap/fonts/glyphicons-halflings-regular.ttf",
	)
}

// css_bootstrap_fonts_glyphicons_halflings_regular_woff reads file data from disk. It returns an error on failure.
func css_bootstrap_fonts_glyphicons_halflings_regular_woff() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/css/bootstrap/fonts/glyphicons-halflings-regular.woff",
		"css/bootstrap/fonts/glyphicons-halflings-regular.woff",
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

// js_ui_bootstrap_custom_tpls_0_10_0_min_js reads file data from disk. It returns an error on failure.
func js_ui_bootstrap_custom_tpls_0_10_0_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/ui-bootstrap-custom-tpls-0.10.0.min.js",
		"js/ui-bootstrap-custom-tpls-0.10.0.min.js",
	)
}

// js_ui_bootstrap_tpls_0_11_0_js reads file data from disk. It returns an error on failure.
func js_ui_bootstrap_tpls_0_11_0_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/ui-bootstrap-tpls-0.11.0.js",
		"js/ui-bootstrap-tpls-0.11.0.js",
	)
}

// js_ui_bootstrap_tpls_0_11_0_min_js reads file data from disk. It returns an error on failure.
func js_ui_bootstrap_tpls_0_11_0_min_js() ([]byte, error) {
	return bindata_read(
		"/Users/byron/Documents/dev/go-lib/src/github.com/Byron/godi/web/client/app/js/ui-bootstrap-tpls-0.11.0.min.js",
		"js/ui-bootstrap-tpls-0.11.0.min.js",
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
	"bower_components/angular/angular-csp.css": bower_components_angular_angular_csp_css,
	"bower_components/angular/angular.js": bower_components_angular_angular_js,
	"bower_components/angular/angular.min.js": bower_components_angular_angular_min_js,
	"bower_components/angular/angular.min.js.gzip": bower_components_angular_angular_min_js_gzip,
	"bower_components/angular/angular.min.js.map": bower_components_angular_angular_min_js_map,
	"bower_components/angular/bower.json": bower_components_angular_bower_json,
	"bower_components/angular/README.md": bower_components_angular_readme_md,
	"bower_components/angular-loader/angular-loader.js": bower_components_angular_loader_angular_loader_js,
	"bower_components/angular-loader/angular-loader.min.js": bower_components_angular_loader_angular_loader_min_js,
	"bower_components/angular-loader/angular-loader.min.js.map": bower_components_angular_loader_angular_loader_min_js_map,
	"bower_components/angular-loader/bower.json": bower_components_angular_loader_bower_json,
	"bower_components/angular-loader/README.md": bower_components_angular_loader_readme_md,
	"bower_components/angular-mocks/angular-mocks.js": bower_components_angular_mocks_angular_mocks_js,
	"bower_components/angular-mocks/bower.json": bower_components_angular_mocks_bower_json,
	"bower_components/angular-mocks/README.md": bower_components_angular_mocks_readme_md,
	"bower_components/angular-resource/angular-resource.js": bower_components_angular_resource_angular_resource_js,
	"bower_components/angular-resource/angular-resource.min.js": bower_components_angular_resource_angular_resource_min_js,
	"bower_components/angular-resource/angular-resource.min.js.map": bower_components_angular_resource_angular_resource_min_js_map,
	"bower_components/angular-resource/bower.json": bower_components_angular_resource_bower_json,
	"bower_components/angular-resource/README.md": bower_components_angular_resource_readme_md,
	"bower_components/angular-route/angular-route.js": bower_components_angular_route_angular_route_js,
	"bower_components/angular-route/angular-route.min.js": bower_components_angular_route_angular_route_min_js,
	"bower_components/angular-route/angular-route.min.js.map": bower_components_angular_route_angular_route_min_js_map,
	"bower_components/angular-route/bower.json": bower_components_angular_route_bower_json,
	"bower_components/angular-route/README.md": bower_components_angular_route_readme_md,
	"bower_components/html5-boilerplate/404.html": bower_components_html5_boilerplate_404_html,
	"bower_components/html5-boilerplate/apple-touch-icon-precomposed.png": bower_components_html5_boilerplate_apple_touch_icon_precomposed_png,
	"bower_components/html5-boilerplate/CHANGELOG.md": bower_components_html5_boilerplate_changelog_md,
	"bower_components/html5-boilerplate/CONTRIBUTING.md": bower_components_html5_boilerplate_contributing_md,
	"bower_components/html5-boilerplate/crossdomain.xml": bower_components_html5_boilerplate_crossdomain_xml,
	"bower_components/html5-boilerplate/css/main.css": bower_components_html5_boilerplate_css_main_css,
	"bower_components/html5-boilerplate/css/normalize.css": bower_components_html5_boilerplate_css_normalize_css,
	"bower_components/html5-boilerplate/doc/crossdomain.md": bower_components_html5_boilerplate_doc_crossdomain_md,
	"bower_components/html5-boilerplate/doc/css.md": bower_components_html5_boilerplate_doc_css_md,
	"bower_components/html5-boilerplate/doc/extend.md": bower_components_html5_boilerplate_doc_extend_md,
	"bower_components/html5-boilerplate/doc/faq.md": bower_components_html5_boilerplate_doc_faq_md,
	"bower_components/html5-boilerplate/doc/html.md": bower_components_html5_boilerplate_doc_html_md,
	"bower_components/html5-boilerplate/doc/js.md": bower_components_html5_boilerplate_doc_js_md,
	"bower_components/html5-boilerplate/doc/misc.md": bower_components_html5_boilerplate_doc_misc_md,
	"bower_components/html5-boilerplate/doc/TOC.md": bower_components_html5_boilerplate_doc_toc_md,
	"bower_components/html5-boilerplate/doc/usage.md": bower_components_html5_boilerplate_doc_usage_md,
	"bower_components/html5-boilerplate/favicon.ico": bower_components_html5_boilerplate_favicon_ico,
	"bower_components/html5-boilerplate/humans.txt": bower_components_html5_boilerplate_humans_txt,
	"bower_components/html5-boilerplate/index.html": bower_components_html5_boilerplate_index_html,
	"bower_components/html5-boilerplate/js/main.js": bower_components_html5_boilerplate_js_main_js,
	"bower_components/html5-boilerplate/js/plugins.js": bower_components_html5_boilerplate_js_plugins_js,
	"bower_components/html5-boilerplate/js/vendor/jquery-1.10.2.min.js": bower_components_html5_boilerplate_js_vendor_jquery_1_10_2_min_js,
	"bower_components/html5-boilerplate/js/vendor/modernizr-2.6.2.min.js": bower_components_html5_boilerplate_js_vendor_modernizr_2_6_2_min_js,
	"bower_components/html5-boilerplate/LICENSE.md": bower_components_html5_boilerplate_license_md,
	"bower_components/html5-boilerplate/README.md": bower_components_html5_boilerplate_readme_md,
	"bower_components/html5-boilerplate/robots.txt": bower_components_html5_boilerplate_robots_txt,
	"bower_components/jquery/bower.json": bower_components_jquery_bower_json,
	"bower_components/jquery/dist/jquery.js": bower_components_jquery_dist_jquery_js,
	"bower_components/jquery/dist/jquery.min.js": bower_components_jquery_dist_jquery_min_js,
	"bower_components/jquery/dist/jquery.min.map": bower_components_jquery_dist_jquery_min_map,
	"bower_components/jquery/MIT-LICENSE.txt": bower_components_jquery_mit_license_txt,
	"bower_components/jquery/src/ajax/jsonp.js": bower_components_jquery_src_ajax_jsonp_js,
	"bower_components/jquery/src/ajax/load.js": bower_components_jquery_src_ajax_load_js,
	"bower_components/jquery/src/ajax/parseJSON.js": bower_components_jquery_src_ajax_parsejson_js,
	"bower_components/jquery/src/ajax/parseXML.js": bower_components_jquery_src_ajax_parsexml_js,
	"bower_components/jquery/src/ajax/script.js": bower_components_jquery_src_ajax_script_js,
	"bower_components/jquery/src/ajax/var/nonce.js": bower_components_jquery_src_ajax_var_nonce_js,
	"bower_components/jquery/src/ajax/var/rquery.js": bower_components_jquery_src_ajax_var_rquery_js,
	"bower_components/jquery/src/ajax/xhr.js": bower_components_jquery_src_ajax_xhr_js,
	"bower_components/jquery/src/ajax.js": bower_components_jquery_src_ajax_js,
	"bower_components/jquery/src/attributes/attr.js": bower_components_jquery_src_attributes_attr_js,
	"bower_components/jquery/src/attributes/classes.js": bower_components_jquery_src_attributes_classes_js,
	"bower_components/jquery/src/attributes/prop.js": bower_components_jquery_src_attributes_prop_js,
	"bower_components/jquery/src/attributes/support.js": bower_components_jquery_src_attributes_support_js,
	"bower_components/jquery/src/attributes/val.js": bower_components_jquery_src_attributes_val_js,
	"bower_components/jquery/src/attributes.js": bower_components_jquery_src_attributes_js,
	"bower_components/jquery/src/callbacks.js": bower_components_jquery_src_callbacks_js,
	"bower_components/jquery/src/core/access.js": bower_components_jquery_src_core_access_js,
	"bower_components/jquery/src/core/init.js": bower_components_jquery_src_core_init_js,
	"bower_components/jquery/src/core/parseHTML.js": bower_components_jquery_src_core_parsehtml_js,
	"bower_components/jquery/src/core/ready.js": bower_components_jquery_src_core_ready_js,
	"bower_components/jquery/src/core/var/rsingleTag.js": bower_components_jquery_src_core_var_rsingletag_js,
	"bower_components/jquery/src/core.js": bower_components_jquery_src_core_js,
	"bower_components/jquery/src/css/addGetHookIf.js": bower_components_jquery_src_css_addgethookif_js,
	"bower_components/jquery/src/css/curCSS.js": bower_components_jquery_src_css_curcss_js,
	"bower_components/jquery/src/css/defaultDisplay.js": bower_components_jquery_src_css_defaultdisplay_js,
	"bower_components/jquery/src/css/hiddenVisibleSelectors.js": bower_components_jquery_src_css_hiddenvisibleselectors_js,
	"bower_components/jquery/src/css/support.js": bower_components_jquery_src_css_support_js,
	"bower_components/jquery/src/css/swap.js": bower_components_jquery_src_css_swap_js,
	"bower_components/jquery/src/css/var/cssExpand.js": bower_components_jquery_src_css_var_cssexpand_js,
	"bower_components/jquery/src/css/var/getStyles.js": bower_components_jquery_src_css_var_getstyles_js,
	"bower_components/jquery/src/css/var/isHidden.js": bower_components_jquery_src_css_var_ishidden_js,
	"bower_components/jquery/src/css/var/rmargin.js": bower_components_jquery_src_css_var_rmargin_js,
	"bower_components/jquery/src/css/var/rnumnonpx.js": bower_components_jquery_src_css_var_rnumnonpx_js,
	"bower_components/jquery/src/css.js": bower_components_jquery_src_css_js,
	"bower_components/jquery/src/data/accepts.js": bower_components_jquery_src_data_accepts_js,
	"bower_components/jquery/src/data/Data.js": bower_components_jquery_src_data_data_js,
	"bower_components/jquery/src/data/var/data_priv.js": bower_components_jquery_src_data_var_data_priv_js,
	"bower_components/jquery/src/data/var/data_user.js": bower_components_jquery_src_data_var_data_user_js,
	"bower_components/jquery/src/data.js": bower_components_jquery_src_data_js,
	"bower_components/jquery/src/deferred.js": bower_components_jquery_src_deferred_js,
	"bower_components/jquery/src/deprecated.js": bower_components_jquery_src_deprecated_js,
	"bower_components/jquery/src/dimensions.js": bower_components_jquery_src_dimensions_js,
	"bower_components/jquery/src/effects/animatedSelector.js": bower_components_jquery_src_effects_animatedselector_js,
	"bower_components/jquery/src/effects/Tween.js": bower_components_jquery_src_effects_tween_js,
	"bower_components/jquery/src/effects.js": bower_components_jquery_src_effects_js,
	"bower_components/jquery/src/event/alias.js": bower_components_jquery_src_event_alias_js,
	"bower_components/jquery/src/event/support.js": bower_components_jquery_src_event_support_js,
	"bower_components/jquery/src/event.js": bower_components_jquery_src_event_js,
	"bower_components/jquery/src/exports/amd.js": bower_components_jquery_src_exports_amd_js,
	"bower_components/jquery/src/exports/global.js": bower_components_jquery_src_exports_global_js,
	"bower_components/jquery/src/intro.js": bower_components_jquery_src_intro_js,
	"bower_components/jquery/src/jquery.js": bower_components_jquery_src_jquery_js,
	"bower_components/jquery/src/manipulation/_evalUrl.js": bower_components_jquery_src_manipulation_evalurl_js,
	"bower_components/jquery/src/manipulation/support.js": bower_components_jquery_src_manipulation_support_js,
	"bower_components/jquery/src/manipulation/var/rcheckableType.js": bower_components_jquery_src_manipulation_var_rcheckabletype_js,
	"bower_components/jquery/src/manipulation.js": bower_components_jquery_src_manipulation_js,
	"bower_components/jquery/src/offset.js": bower_components_jquery_src_offset_js,
	"bower_components/jquery/src/outro.js": bower_components_jquery_src_outro_js,
	"bower_components/jquery/src/queue/delay.js": bower_components_jquery_src_queue_delay_js,
	"bower_components/jquery/src/queue.js": bower_components_jquery_src_queue_js,
	"bower_components/jquery/src/selector-native.js": bower_components_jquery_src_selector_native_js,
	"bower_components/jquery/src/selector-sizzle.js": bower_components_jquery_src_selector_sizzle_js,
	"bower_components/jquery/src/selector.js": bower_components_jquery_src_selector_js,
	"bower_components/jquery/src/serialize.js": bower_components_jquery_src_serialize_js,
	"bower_components/jquery/src/sizzle/dist/sizzle.js": bower_components_jquery_src_sizzle_dist_sizzle_js,
	"bower_components/jquery/src/sizzle/dist/sizzle.min.js": bower_components_jquery_src_sizzle_dist_sizzle_min_js,
	"bower_components/jquery/src/sizzle/dist/sizzle.min.map": bower_components_jquery_src_sizzle_dist_sizzle_min_map,
	"bower_components/jquery/src/traversing/findFilter.js": bower_components_jquery_src_traversing_findfilter_js,
	"bower_components/jquery/src/traversing/var/rneedsContext.js": bower_components_jquery_src_traversing_var_rneedscontext_js,
	"bower_components/jquery/src/traversing.js": bower_components_jquery_src_traversing_js,
	"bower_components/jquery/src/var/arr.js": bower_components_jquery_src_var_arr_js,
	"bower_components/jquery/src/var/class2type.js": bower_components_jquery_src_var_class2type_js,
	"bower_components/jquery/src/var/concat.js": bower_components_jquery_src_var_concat_js,
	"bower_components/jquery/src/var/hasOwn.js": bower_components_jquery_src_var_hasown_js,
	"bower_components/jquery/src/var/indexOf.js": bower_components_jquery_src_var_indexof_js,
	"bower_components/jquery/src/var/pnum.js": bower_components_jquery_src_var_pnum_js,
	"bower_components/jquery/src/var/push.js": bower_components_jquery_src_var_push_js,
	"bower_components/jquery/src/var/rnotwhite.js": bower_components_jquery_src_var_rnotwhite_js,
	"bower_components/jquery/src/var/slice.js": bower_components_jquery_src_var_slice_js,
	"bower_components/jquery/src/var/strundefined.js": bower_components_jquery_src_var_strundefined_js,
	"bower_components/jquery/src/var/support.js": bower_components_jquery_src_var_support_js,
	"bower_components/jquery/src/var/toString.js": bower_components_jquery_src_var_tostring_js,
	"bower_components/jquery/src/wrap.js": bower_components_jquery_src_wrap_js,
	"css/app.css": css_app_css,
	"css/bootstrap/css/bootstrap-theme.css": css_bootstrap_css_bootstrap_theme_css,
	"css/bootstrap/css/bootstrap-theme.css.map": css_bootstrap_css_bootstrap_theme_css_map,
	"css/bootstrap/css/bootstrap-theme.min.css": css_bootstrap_css_bootstrap_theme_min_css,
	"css/bootstrap/css/bootstrap.css": css_bootstrap_css_bootstrap_css,
	"css/bootstrap/css/bootstrap.css.map": css_bootstrap_css_bootstrap_css_map,
	"css/bootstrap/css/bootstrap.min.css": css_bootstrap_css_bootstrap_min_css,
	"css/bootstrap/fonts/glyphicons-halflings-regular.eot": css_bootstrap_fonts_glyphicons_halflings_regular_eot,
	"css/bootstrap/fonts/glyphicons-halflings-regular.svg": css_bootstrap_fonts_glyphicons_halflings_regular_svg,
	"css/bootstrap/fonts/glyphicons-halflings-regular.ttf": css_bootstrap_fonts_glyphicons_halflings_regular_ttf,
	"css/bootstrap/fonts/glyphicons-halflings-regular.woff": css_bootstrap_fonts_glyphicons_halflings_regular_woff,
	"index-async.html": index_async_html,
	"index.html": index_html,
	"js/app.js": js_app_js,
	"js/controllers.js": js_controllers_js,
	"js/directives.js": js_directives_js,
	"js/filters.js": js_filters_js,
	"js/services.js": js_services_js,
	"js/ui-bootstrap-custom-tpls-0.10.0.min.js": js_ui_bootstrap_custom_tpls_0_10_0_min_js,
	"js/ui-bootstrap-tpls-0.11.0.js": js_ui_bootstrap_tpls_0_11_0_js,
	"js/ui-bootstrap-tpls-0.11.0.min.js": js_ui_bootstrap_tpls_0_11_0_min_js,
	"npm-debug.log": npm_debug_log,
	"partials/partial1.html": partials_partial1_html,
	"partials/partial2.html": partials_partial2_html,
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
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"index-async.html": &_bintree_t{index_async_html, map[string]*_bintree_t{
	}},
	"index.html": &_bintree_t{index_html, map[string]*_bintree_t{
	}},
	"js": &_bintree_t{nil, map[string]*_bintree_t{
		"filters.js": &_bintree_t{js_filters_js, map[string]*_bintree_t{
		}},
		"services.js": &_bintree_t{js_services_js, map[string]*_bintree_t{
		}},
		"ui-bootstrap-custom-tpls-0.10.0.min.js": &_bintree_t{js_ui_bootstrap_custom_tpls_0_10_0_min_js, map[string]*_bintree_t{
		}},
		"ui-bootstrap-tpls-0.11.0.js": &_bintree_t{js_ui_bootstrap_tpls_0_11_0_js, map[string]*_bintree_t{
		}},
		"ui-bootstrap-tpls-0.11.0.min.js": &_bintree_t{js_ui_bootstrap_tpls_0_11_0_min_js, map[string]*_bintree_t{
		}},
		"app.js": &_bintree_t{js_app_js, map[string]*_bintree_t{
		}},
		"controllers.js": &_bintree_t{js_controllers_js, map[string]*_bintree_t{
		}},
		"directives.js": &_bintree_t{js_directives_js, map[string]*_bintree_t{
		}},
	}},
	"npm-debug.log": &_bintree_t{npm_debug_log, map[string]*_bintree_t{
	}},
	"partials": &_bintree_t{nil, map[string]*_bintree_t{
		"partial2.html": &_bintree_t{partials_partial2_html, map[string]*_bintree_t{
		}},
		"partial1.html": &_bintree_t{partials_partial1_html, map[string]*_bintree_t{
		}},
	}},
	"bower_components": &_bintree_t{nil, map[string]*_bintree_t{
		"angular-mocks": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-mocks.js": &_bintree_t{bower_components_angular_mocks_angular_mocks_js, map[string]*_bintree_t{
			}},
			"bower.json": &_bintree_t{bower_components_angular_mocks_bower_json, map[string]*_bintree_t{
			}},
			"README.md": &_bintree_t{bower_components_angular_mocks_readme_md, map[string]*_bintree_t{
			}},
		}},
		"angular-resource": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-resource.js": &_bintree_t{bower_components_angular_resource_angular_resource_js, map[string]*_bintree_t{
			}},
			"angular-resource.min.js": &_bintree_t{bower_components_angular_resource_angular_resource_min_js, map[string]*_bintree_t{
			}},
			"angular-resource.min.js.map": &_bintree_t{bower_components_angular_resource_angular_resource_min_js_map, map[string]*_bintree_t{
			}},
			"bower.json": &_bintree_t{bower_components_angular_resource_bower_json, map[string]*_bintree_t{
			}},
			"README.md": &_bintree_t{bower_components_angular_resource_readme_md, map[string]*_bintree_t{
			}},
		}},
		"angular-route": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-route.min.js": &_bintree_t{bower_components_angular_route_angular_route_min_js, map[string]*_bintree_t{
			}},
			"angular-route.min.js.map": &_bintree_t{bower_components_angular_route_angular_route_min_js_map, map[string]*_bintree_t{
			}},
			"bower.json": &_bintree_t{bower_components_angular_route_bower_json, map[string]*_bintree_t{
			}},
			"README.md": &_bintree_t{bower_components_angular_route_readme_md, map[string]*_bintree_t{
			}},
			"angular-route.js": &_bintree_t{bower_components_angular_route_angular_route_js, map[string]*_bintree_t{
			}},
		}},
		"html5-boilerplate": &_bintree_t{nil, map[string]*_bintree_t{
			"404.html": &_bintree_t{bower_components_html5_boilerplate_404_html, map[string]*_bintree_t{
			}},
			"apple-touch-icon-precomposed.png": &_bintree_t{bower_components_html5_boilerplate_apple_touch_icon_precomposed_png, map[string]*_bintree_t{
			}},
			"README.md": &_bintree_t{bower_components_html5_boilerplate_readme_md, map[string]*_bintree_t{
			}},
			"favicon.ico": &_bintree_t{bower_components_html5_boilerplate_favicon_ico, map[string]*_bintree_t{
			}},
			"js": &_bintree_t{nil, map[string]*_bintree_t{
				"vendor": &_bintree_t{nil, map[string]*_bintree_t{
					"jquery-1.10.2.min.js": &_bintree_t{bower_components_html5_boilerplate_js_vendor_jquery_1_10_2_min_js, map[string]*_bintree_t{
					}},
					"modernizr-2.6.2.min.js": &_bintree_t{bower_components_html5_boilerplate_js_vendor_modernizr_2_6_2_min_js, map[string]*_bintree_t{
					}},
				}},
				"main.js": &_bintree_t{bower_components_html5_boilerplate_js_main_js, map[string]*_bintree_t{
				}},
				"plugins.js": &_bintree_t{bower_components_html5_boilerplate_js_plugins_js, map[string]*_bintree_t{
				}},
			}},
			"robots.txt": &_bintree_t{bower_components_html5_boilerplate_robots_txt, map[string]*_bintree_t{
			}},
			"CONTRIBUTING.md": &_bintree_t{bower_components_html5_boilerplate_contributing_md, map[string]*_bintree_t{
			}},
			"crossdomain.xml": &_bintree_t{bower_components_html5_boilerplate_crossdomain_xml, map[string]*_bintree_t{
			}},
			"doc": &_bintree_t{nil, map[string]*_bintree_t{
				"faq.md": &_bintree_t{bower_components_html5_boilerplate_doc_faq_md, map[string]*_bintree_t{
				}},
				"html.md": &_bintree_t{bower_components_html5_boilerplate_doc_html_md, map[string]*_bintree_t{
				}},
				"js.md": &_bintree_t{bower_components_html5_boilerplate_doc_js_md, map[string]*_bintree_t{
				}},
				"crossdomain.md": &_bintree_t{bower_components_html5_boilerplate_doc_crossdomain_md, map[string]*_bintree_t{
				}},
				"extend.md": &_bintree_t{bower_components_html5_boilerplate_doc_extend_md, map[string]*_bintree_t{
				}},
				"misc.md": &_bintree_t{bower_components_html5_boilerplate_doc_misc_md, map[string]*_bintree_t{
				}},
				"TOC.md": &_bintree_t{bower_components_html5_boilerplate_doc_toc_md, map[string]*_bintree_t{
				}},
				"usage.md": &_bintree_t{bower_components_html5_boilerplate_doc_usage_md, map[string]*_bintree_t{
				}},
				"css.md": &_bintree_t{bower_components_html5_boilerplate_doc_css_md, map[string]*_bintree_t{
				}},
			}},
			"index.html": &_bintree_t{bower_components_html5_boilerplate_index_html, map[string]*_bintree_t{
			}},
			"CHANGELOG.md": &_bintree_t{bower_components_html5_boilerplate_changelog_md, map[string]*_bintree_t{
			}},
			"css": &_bintree_t{nil, map[string]*_bintree_t{
				"main.css": &_bintree_t{bower_components_html5_boilerplate_css_main_css, map[string]*_bintree_t{
				}},
				"normalize.css": &_bintree_t{bower_components_html5_boilerplate_css_normalize_css, map[string]*_bintree_t{
				}},
			}},
			"humans.txt": &_bintree_t{bower_components_html5_boilerplate_humans_txt, map[string]*_bintree_t{
			}},
			"LICENSE.md": &_bintree_t{bower_components_html5_boilerplate_license_md, map[string]*_bintree_t{
			}},
		}},
		"jquery": &_bintree_t{nil, map[string]*_bintree_t{
			"src": &_bintree_t{nil, map[string]*_bintree_t{
				"css": &_bintree_t{nil, map[string]*_bintree_t{
					"defaultDisplay.js": &_bintree_t{bower_components_jquery_src_css_defaultdisplay_js, map[string]*_bintree_t{
					}},
					"hiddenVisibleSelectors.js": &_bintree_t{bower_components_jquery_src_css_hiddenvisibleselectors_js, map[string]*_bintree_t{
					}},
					"support.js": &_bintree_t{bower_components_jquery_src_css_support_js, map[string]*_bintree_t{
					}},
					"swap.js": &_bintree_t{bower_components_jquery_src_css_swap_js, map[string]*_bintree_t{
					}},
					"var": &_bintree_t{nil, map[string]*_bintree_t{
						"rmargin.js": &_bintree_t{bower_components_jquery_src_css_var_rmargin_js, map[string]*_bintree_t{
						}},
						"rnumnonpx.js": &_bintree_t{bower_components_jquery_src_css_var_rnumnonpx_js, map[string]*_bintree_t{
						}},
						"cssExpand.js": &_bintree_t{bower_components_jquery_src_css_var_cssexpand_js, map[string]*_bintree_t{
						}},
						"getStyles.js": &_bintree_t{bower_components_jquery_src_css_var_getstyles_js, map[string]*_bintree_t{
						}},
						"isHidden.js": &_bintree_t{bower_components_jquery_src_css_var_ishidden_js, map[string]*_bintree_t{
						}},
					}},
					"addGetHookIf.js": &_bintree_t{bower_components_jquery_src_css_addgethookif_js, map[string]*_bintree_t{
					}},
					"curCSS.js": &_bintree_t{bower_components_jquery_src_css_curcss_js, map[string]*_bintree_t{
					}},
				}},
				"data": &_bintree_t{nil, map[string]*_bintree_t{
					"Data.js": &_bintree_t{bower_components_jquery_src_data_data_js, map[string]*_bintree_t{
					}},
					"var": &_bintree_t{nil, map[string]*_bintree_t{
						"data_priv.js": &_bintree_t{bower_components_jquery_src_data_var_data_priv_js, map[string]*_bintree_t{
						}},
						"data_user.js": &_bintree_t{bower_components_jquery_src_data_var_data_user_js, map[string]*_bintree_t{
						}},
					}},
					"accepts.js": &_bintree_t{bower_components_jquery_src_data_accepts_js, map[string]*_bintree_t{
					}},
				}},
				"data.js": &_bintree_t{bower_components_jquery_src_data_js, map[string]*_bintree_t{
				}},
				"intro.js": &_bintree_t{bower_components_jquery_src_intro_js, map[string]*_bintree_t{
				}},
				"queue.js": &_bintree_t{bower_components_jquery_src_queue_js, map[string]*_bintree_t{
				}},
				"selector-sizzle.js": &_bintree_t{bower_components_jquery_src_selector_sizzle_js, map[string]*_bintree_t{
				}},
				"selector.js": &_bintree_t{bower_components_jquery_src_selector_js, map[string]*_bintree_t{
				}},
				"core": &_bintree_t{nil, map[string]*_bintree_t{
					"ready.js": &_bintree_t{bower_components_jquery_src_core_ready_js, map[string]*_bintree_t{
					}},
					"var": &_bintree_t{nil, map[string]*_bintree_t{
						"rsingleTag.js": &_bintree_t{bower_components_jquery_src_core_var_rsingletag_js, map[string]*_bintree_t{
						}},
					}},
					"access.js": &_bintree_t{bower_components_jquery_src_core_access_js, map[string]*_bintree_t{
					}},
					"init.js": &_bintree_t{bower_components_jquery_src_core_init_js, map[string]*_bintree_t{
					}},
					"parseHTML.js": &_bintree_t{bower_components_jquery_src_core_parsehtml_js, map[string]*_bintree_t{
					}},
				}},
				"serialize.js": &_bintree_t{bower_components_jquery_src_serialize_js, map[string]*_bintree_t{
				}},
				"sizzle": &_bintree_t{nil, map[string]*_bintree_t{
					"dist": &_bintree_t{nil, map[string]*_bintree_t{
						"sizzle.js": &_bintree_t{bower_components_jquery_src_sizzle_dist_sizzle_js, map[string]*_bintree_t{
						}},
						"sizzle.min.js": &_bintree_t{bower_components_jquery_src_sizzle_dist_sizzle_min_js, map[string]*_bintree_t{
						}},
						"sizzle.min.map": &_bintree_t{bower_components_jquery_src_sizzle_dist_sizzle_min_map, map[string]*_bintree_t{
						}},
					}},
				}},
				"ajax": &_bintree_t{nil, map[string]*_bintree_t{
					"parseJSON.js": &_bintree_t{bower_components_jquery_src_ajax_parsejson_js, map[string]*_bintree_t{
					}},
					"parseXML.js": &_bintree_t{bower_components_jquery_src_ajax_parsexml_js, map[string]*_bintree_t{
					}},
					"script.js": &_bintree_t{bower_components_jquery_src_ajax_script_js, map[string]*_bintree_t{
					}},
					"var": &_bintree_t{nil, map[string]*_bintree_t{
						"nonce.js": &_bintree_t{bower_components_jquery_src_ajax_var_nonce_js, map[string]*_bintree_t{
						}},
						"rquery.js": &_bintree_t{bower_components_jquery_src_ajax_var_rquery_js, map[string]*_bintree_t{
						}},
					}},
					"xhr.js": &_bintree_t{bower_components_jquery_src_ajax_xhr_js, map[string]*_bintree_t{
					}},
					"jsonp.js": &_bintree_t{bower_components_jquery_src_ajax_jsonp_js, map[string]*_bintree_t{
					}},
					"load.js": &_bintree_t{bower_components_jquery_src_ajax_load_js, map[string]*_bintree_t{
					}},
				}},
				"attributes.js": &_bintree_t{bower_components_jquery_src_attributes_js, map[string]*_bintree_t{
				}},
				"attributes": &_bintree_t{nil, map[string]*_bintree_t{
					"classes.js": &_bintree_t{bower_components_jquery_src_attributes_classes_js, map[string]*_bintree_t{
					}},
					"prop.js": &_bintree_t{bower_components_jquery_src_attributes_prop_js, map[string]*_bintree_t{
					}},
					"support.js": &_bintree_t{bower_components_jquery_src_attributes_support_js, map[string]*_bintree_t{
					}},
					"val.js": &_bintree_t{bower_components_jquery_src_attributes_val_js, map[string]*_bintree_t{
					}},
					"attr.js": &_bintree_t{bower_components_jquery_src_attributes_attr_js, map[string]*_bintree_t{
					}},
				}},
				"manipulation": &_bintree_t{nil, map[string]*_bintree_t{
					"var": &_bintree_t{nil, map[string]*_bintree_t{
						"rcheckableType.js": &_bintree_t{bower_components_jquery_src_manipulation_var_rcheckabletype_js, map[string]*_bintree_t{
						}},
					}},
					"_evalUrl.js": &_bintree_t{bower_components_jquery_src_manipulation_evalurl_js, map[string]*_bintree_t{
					}},
					"support.js": &_bintree_t{bower_components_jquery_src_manipulation_support_js, map[string]*_bintree_t{
					}},
				}},
				"var": &_bintree_t{nil, map[string]*_bintree_t{
					"class2type.js": &_bintree_t{bower_components_jquery_src_var_class2type_js, map[string]*_bintree_t{
					}},
					"hasOwn.js": &_bintree_t{bower_components_jquery_src_var_hasown_js, map[string]*_bintree_t{
					}},
					"indexOf.js": &_bintree_t{bower_components_jquery_src_var_indexof_js, map[string]*_bintree_t{
					}},
					"pnum.js": &_bintree_t{bower_components_jquery_src_var_pnum_js, map[string]*_bintree_t{
					}},
					"slice.js": &_bintree_t{bower_components_jquery_src_var_slice_js, map[string]*_bintree_t{
					}},
					"arr.js": &_bintree_t{bower_components_jquery_src_var_arr_js, map[string]*_bintree_t{
					}},
					"push.js": &_bintree_t{bower_components_jquery_src_var_push_js, map[string]*_bintree_t{
					}},
					"rnotwhite.js": &_bintree_t{bower_components_jquery_src_var_rnotwhite_js, map[string]*_bintree_t{
					}},
					"strundefined.js": &_bintree_t{bower_components_jquery_src_var_strundefined_js, map[string]*_bintree_t{
					}},
					"support.js": &_bintree_t{bower_components_jquery_src_var_support_js, map[string]*_bintree_t{
					}},
					"toString.js": &_bintree_t{bower_components_jquery_src_var_tostring_js, map[string]*_bintree_t{
					}},
					"concat.js": &_bintree_t{bower_components_jquery_src_var_concat_js, map[string]*_bintree_t{
					}},
				}},
				"css.js": &_bintree_t{bower_components_jquery_src_css_js, map[string]*_bintree_t{
				}},
				"effects": &_bintree_t{nil, map[string]*_bintree_t{
					"animatedSelector.js": &_bintree_t{bower_components_jquery_src_effects_animatedselector_js, map[string]*_bintree_t{
					}},
					"Tween.js": &_bintree_t{bower_components_jquery_src_effects_tween_js, map[string]*_bintree_t{
					}},
				}},
				"event": &_bintree_t{nil, map[string]*_bintree_t{
					"alias.js": &_bintree_t{bower_components_jquery_src_event_alias_js, map[string]*_bintree_t{
					}},
					"support.js": &_bintree_t{bower_components_jquery_src_event_support_js, map[string]*_bintree_t{
					}},
				}},
				"callbacks.js": &_bintree_t{bower_components_jquery_src_callbacks_js, map[string]*_bintree_t{
				}},
				"core.js": &_bintree_t{bower_components_jquery_src_core_js, map[string]*_bintree_t{
				}},
				"deprecated.js": &_bintree_t{bower_components_jquery_src_deprecated_js, map[string]*_bintree_t{
				}},
				"offset.js": &_bintree_t{bower_components_jquery_src_offset_js, map[string]*_bintree_t{
				}},
				"queue": &_bintree_t{nil, map[string]*_bintree_t{
					"delay.js": &_bintree_t{bower_components_jquery_src_queue_delay_js, map[string]*_bintree_t{
					}},
				}},
				"traversing.js": &_bintree_t{bower_components_jquery_src_traversing_js, map[string]*_bintree_t{
				}},
				"wrap.js": &_bintree_t{bower_components_jquery_src_wrap_js, map[string]*_bintree_t{
				}},
				"ajax.js": &_bintree_t{bower_components_jquery_src_ajax_js, map[string]*_bintree_t{
				}},
				"effects.js": &_bintree_t{bower_components_jquery_src_effects_js, map[string]*_bintree_t{
				}},
				"exports": &_bintree_t{nil, map[string]*_bintree_t{
					"global.js": &_bintree_t{bower_components_jquery_src_exports_global_js, map[string]*_bintree_t{
					}},
					"amd.js": &_bintree_t{bower_components_jquery_src_exports_amd_js, map[string]*_bintree_t{
					}},
				}},
				"outro.js": &_bintree_t{bower_components_jquery_src_outro_js, map[string]*_bintree_t{
				}},
				"deferred.js": &_bintree_t{bower_components_jquery_src_deferred_js, map[string]*_bintree_t{
				}},
				"event.js": &_bintree_t{bower_components_jquery_src_event_js, map[string]*_bintree_t{
				}},
				"jquery.js": &_bintree_t{bower_components_jquery_src_jquery_js, map[string]*_bintree_t{
				}},
				"manipulation.js": &_bintree_t{bower_components_jquery_src_manipulation_js, map[string]*_bintree_t{
				}},
				"selector-native.js": &_bintree_t{bower_components_jquery_src_selector_native_js, map[string]*_bintree_t{
				}},
				"traversing": &_bintree_t{nil, map[string]*_bintree_t{
					"findFilter.js": &_bintree_t{bower_components_jquery_src_traversing_findfilter_js, map[string]*_bintree_t{
					}},
					"var": &_bintree_t{nil, map[string]*_bintree_t{
						"rneedsContext.js": &_bintree_t{bower_components_jquery_src_traversing_var_rneedscontext_js, map[string]*_bintree_t{
						}},
					}},
				}},
				"dimensions.js": &_bintree_t{bower_components_jquery_src_dimensions_js, map[string]*_bintree_t{
				}},
			}},
			"bower.json": &_bintree_t{bower_components_jquery_bower_json, map[string]*_bintree_t{
			}},
			"dist": &_bintree_t{nil, map[string]*_bintree_t{
				"jquery.js": &_bintree_t{bower_components_jquery_dist_jquery_js, map[string]*_bintree_t{
				}},
				"jquery.min.js": &_bintree_t{bower_components_jquery_dist_jquery_min_js, map[string]*_bintree_t{
				}},
				"jquery.min.map": &_bintree_t{bower_components_jquery_dist_jquery_min_map, map[string]*_bintree_t{
				}},
			}},
			"MIT-LICENSE.txt": &_bintree_t{bower_components_jquery_mit_license_txt, map[string]*_bintree_t{
			}},
		}},
		"angular": &_bintree_t{nil, map[string]*_bintree_t{
			"bower.json": &_bintree_t{bower_components_angular_bower_json, map[string]*_bintree_t{
			}},
			"README.md": &_bintree_t{bower_components_angular_readme_md, map[string]*_bintree_t{
			}},
			"angular-csp.css": &_bintree_t{bower_components_angular_angular_csp_css, map[string]*_bintree_t{
			}},
			"angular.js": &_bintree_t{bower_components_angular_angular_js, map[string]*_bintree_t{
			}},
			"angular.min.js": &_bintree_t{bower_components_angular_angular_min_js, map[string]*_bintree_t{
			}},
			"angular.min.js.gzip": &_bintree_t{bower_components_angular_angular_min_js_gzip, map[string]*_bintree_t{
			}},
			"angular.min.js.map": &_bintree_t{bower_components_angular_angular_min_js_map, map[string]*_bintree_t{
			}},
		}},
		"angular-loader": &_bintree_t{nil, map[string]*_bintree_t{
			"angular-loader.min.js.map": &_bintree_t{bower_components_angular_loader_angular_loader_min_js_map, map[string]*_bintree_t{
			}},
			"bower.json": &_bintree_t{bower_components_angular_loader_bower_json, map[string]*_bintree_t{
			}},
			"README.md": &_bintree_t{bower_components_angular_loader_readme_md, map[string]*_bintree_t{
			}},
			"angular-loader.js": &_bintree_t{bower_components_angular_loader_angular_loader_js, map[string]*_bintree_t{
			}},
			"angular-loader.min.js": &_bintree_t{bower_components_angular_loader_angular_loader_min_js, map[string]*_bintree_t{
			}},
		}},
	}},
	"css": &_bintree_t{nil, map[string]*_bintree_t{
		"app.css": &_bintree_t{css_app_css, map[string]*_bintree_t{
		}},
		"bootstrap": &_bintree_t{nil, map[string]*_bintree_t{
			"css": &_bintree_t{nil, map[string]*_bintree_t{
				"bootstrap.css.map": &_bintree_t{css_bootstrap_css_bootstrap_css_map, map[string]*_bintree_t{
				}},
				"bootstrap.min.css": &_bintree_t{css_bootstrap_css_bootstrap_min_css, map[string]*_bintree_t{
				}},
				"bootstrap-theme.css": &_bintree_t{css_bootstrap_css_bootstrap_theme_css, map[string]*_bintree_t{
				}},
				"bootstrap-theme.css.map": &_bintree_t{css_bootstrap_css_bootstrap_theme_css_map, map[string]*_bintree_t{
				}},
				"bootstrap-theme.min.css": &_bintree_t{css_bootstrap_css_bootstrap_theme_min_css, map[string]*_bintree_t{
				}},
				"bootstrap.css": &_bintree_t{css_bootstrap_css_bootstrap_css, map[string]*_bintree_t{
				}},
			}},
			"fonts": &_bintree_t{nil, map[string]*_bintree_t{
				"glyphicons-halflings-regular.eot": &_bintree_t{css_bootstrap_fonts_glyphicons_halflings_regular_eot, map[string]*_bintree_t{
				}},
				"glyphicons-halflings-regular.svg": &_bintree_t{css_bootstrap_fonts_glyphicons_halflings_regular_svg, map[string]*_bintree_t{
				}},
				"glyphicons-halflings-regular.ttf": &_bintree_t{css_bootstrap_fonts_glyphicons_halflings_regular_ttf, map[string]*_bintree_t{
				}},
				"glyphicons-halflings-regular.woff": &_bintree_t{css_bootstrap_fonts_glyphicons_halflings_regular_woff, map[string]*_bintree_t{
				}},
			}},
		}},
	}},
}}
