package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_Configuration_SubConfiguration(t *testing.T) {
	fnName := "Configuration.SubConfiguration()"
	savedState := SaveEnvVarForTesting(appDataVar)
	os.Setenv(appDataVar, SecureAbsolutePathForTesting("."))
	defer func() {
		savedState.RestoreForTesting()
	}()
	if err := CreateDefaultYamlFileForTesting(); err != nil {
		t.Errorf("%s error creating defaults.yaml: %v", fnName, err)
	}
	defer func() {
		DestroyDirectoryForTesting(fnName, "./mp3")
	}()
	testConfiguration, _ := ReadConfigurationFile(NewOutputDeviceForTesting())
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Configuration
		args args
		want *Configuration
	}{
		{name: "no configuration", c: &Configuration{}, args: args{}, want: EmptyConfiguration()},
		{name: "commons", c: testConfiguration, args: args{key: "common"}, want: testConfiguration.cMap["common"]},
		{name: "ls", c: testConfiguration, args: args{key: "ls"}, want: testConfiguration.cMap["ls"]},
		{name: "check", c: testConfiguration, args: args{key: "check"}, want: testConfiguration.cMap["check"]},
		{name: "repair", c: testConfiguration, args: args{key: "repair"}, want: testConfiguration.cMap["repair"]},
		{name: "unknown key", c: testConfiguration, args: args{key: "unknown key"}, want: EmptyConfiguration()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SubConfiguration(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s = %v, want %v", fnName, got, tt.want)
			}
		})
	}
}

func Test_Configuration_BoolDefault(t *testing.T) {
	fnName := "Configuration.BoolDefault()"
	savedState := SaveEnvVarForTesting(appDataVar)
	os.Setenv(appDataVar, SecureAbsolutePathForTesting("."))
	defer func() {
		savedState.RestoreForTesting()
	}()
	if err := CreateDefaultYamlFileForTesting(); err != nil {
		t.Errorf("%s error creating defaults.yaml: %v", fnName, err)
	}
	defer func() {
		DestroyDirectoryForTesting(fnName, "./mp3")
	}()
	testConfiguration, _ := ReadConfigurationFile(NewOutputDeviceForTesting())
	altConfiguration := EmptyConfiguration()
	altConfiguration.sMap["foo"] = "1"
	type args struct {
		key          string
		defaultValue bool
	}
	tests := []struct {
		name  string
		c     *Configuration
		args  args
		wantB bool
	}{
		{
			name:  "string to bool",
			c:     altConfiguration,
			args:  args{key: "foo", defaultValue: false},
			wantB: true,
		},
		{
			name:  "empty configuration default false",
			c:     EmptyConfiguration(),
			args:  args{defaultValue: false},
			wantB: false,
		},
		{
			name:  "empty configuration default true",
			c:     EmptyConfiguration(),
			args:  args{defaultValue: true},
			wantB: true,
		},
		{
			name:  "undefined key default false",
			c:     testConfiguration,
			args:  args{key: "no key", defaultValue: false},
			wantB: false,
		},
		{
			name:  "undefined key default true",
			c:     testConfiguration,
			args:  args{key: "no key", defaultValue: true},
			wantB: true,
		},
		{
			name:  "non-boolean value default false",
			c:     testConfiguration.cMap["common"],
			args:  args{key: "albums", defaultValue: false},
			wantB: false,
		},
		{
			name:  "non-boolean value default true",
			c:     testConfiguration.cMap["common"],
			args:  args{key: "albums", defaultValue: true},
			wantB: true,
		},
		{
			name:  "boolean value default false",
			c:     testConfiguration.cMap["ls"],
			args:  args{key: "includeTracks", defaultValue: false},
			wantB: true,
		},
		{
			name:  "boolean value default true",
			c:     testConfiguration.cMap["ls"],
			args:  args{key: "includeTracks", defaultValue: true},
			wantB: true,
		},
		{
			name:  "boolean (string) value default true",
			c:     testConfiguration.cMap["unused"],
			args:  args{key: "value", defaultValue: true},
			wantB: true,
		},
		{
			name:  "boolean (string) value default false",
			c:     testConfiguration.cMap["unused"],
			args:  args{key: "value", defaultValue: false},
			wantB: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := tt.c.BoolDefault(tt.args.key, tt.args.defaultValue); gotB != tt.wantB {
				t.Errorf("%s = %v, want %v", fnName, gotB, tt.wantB)
			}
		})
	}
}

func Test_Configuration_StringDefault(t *testing.T) {
	fnName := "Configuration.StringDefault()"
	savedState := SaveEnvVarForTesting(appDataVar)
	os.Setenv(appDataVar, SecureAbsolutePathForTesting("."))
	defer func() {
		savedState.RestoreForTesting()
	}()
	if err := CreateDefaultYamlFileForTesting(); err != nil {
		t.Errorf("%s error creating defaults.yaml: %v", fnName, err)
	}
	defer func() {
		DestroyDirectoryForTesting(fnName, "./mp3")
	}()
	testConfiguration, _ := ReadConfigurationFile(NewOutputDeviceForTesting())
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name  string
		c     *Configuration
		args  args
		wantS string
	}{
		{
			name:  "empty configuration",
			c:     EmptyConfiguration(),
			args:  args{defaultValue: "my default value"},
			wantS: "my default value"},
		{
			name:  "undefined key",
			c:     testConfiguration,
			args:  args{key: "no key", defaultValue: "my default value"},
			wantS: "my default value",
		},
		{
			name:  "defined key",
			c:     testConfiguration.cMap["ls"],
			args:  args{key: "sort", defaultValue: "my default value"},
			wantS: "alpha",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := tt.c.StringDefault(tt.args.key, tt.args.defaultValue); gotS != tt.wantS {
				t.Errorf("%s = %v, want %v", fnName, gotS, tt.wantS)
			}
		})
	}
}

func Test_verifyFileExists(t *testing.T) {
	fnName := "verifyFileExists()"
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
		WantedOutput
	}{
		{
			name:         "ordinary success",
			args:         args{path: "./configuration_test.go"},
			wantOk:       true,
			WantedOutput: WantedOutput{},
		},
		{
			name:    "look for dir!",
			args:    args{path: "."},
			wantErr: true,
			WantedOutput: WantedOutput{
				WantErrorOutput: "The configuration file \".\" is a directory.\n",
				WantLogOutput:   "level='error' directory='.' fileName='.' msg='file is a directory'\n",
			},
		},
		{
			name: "non-existent file",
			args: args{path: "./no-such-file.txt"},
			WantedOutput: WantedOutput{
				WantLogOutput: "level='info' directory='.' fileName='no-such-file.txt' msg='file does not exist'\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOutputDeviceForTesting()
			gotOk, err := verifyFileExists(o, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", fnName, err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("%s = %v, want %v", fnName, gotOk, tt.wantOk)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}

func TestReadConfigurationFile(t *testing.T) {
	fnName := "ReadConfigurationFile()"
	savedState := SaveEnvVarForTesting(appDataVar)
	canonicalPath := SecureAbsolutePathForTesting(".")
	mp3Path := SecureAbsolutePathForTesting("mp3")
	os.Setenv(appDataVar, canonicalPath)
	defer func() {
		savedState.RestoreForTesting()
	}()
	if err := CreateDefaultYamlFileForTesting(); err != nil {
		t.Errorf("%s error creating defaults.yaml: %v", fnName, err)
	}
	defer func() {
		DestroyDirectoryForTesting(fnName, "./mp3")
	}()
	badDir := filepath.Join("./mp3", "fake")
	if err := Mkdir(badDir); err != nil {
		t.Errorf("%s error creating fake dir: %v", fnName, err)
	}
	badDir2 := filepath.Join(badDir, AppName)
	if err := Mkdir(badDir2); err != nil {
		t.Errorf("%s error creating fake dir mp3 folder: %v", fnName, err)
	}
	badFile := filepath.Join(badDir2, defaultConfigFileName)
	if err := Mkdir(badFile); err != nil {
		t.Errorf("%s error creating defaults.yaml as a directory: %v", fnName, err)
	}
	yamlAsDir := SecureAbsolutePathForTesting(badFile)
	gibberishDir := filepath.Join(badDir2, AppName)
	if err := Mkdir(gibberishDir); err != nil {
		t.Errorf("%s error creating gibberish folder: %v", fnName, err)
	}
	if err := CreateFileForTestingWithContent(gibberishDir, defaultConfigFileName, []byte("gibberish")); err != nil {
		t.Errorf("%s error creating gibberish defaults.yaml: %v", fnName, err)
	}
	tests := []struct {
		name   string
		state  *SavedEnvVar
		wantC  *Configuration
		wantOk bool
		WantedOutput
	}{
		{
			name:  "good",
			state: &SavedEnvVar{Name: appDataVar, Value: canonicalPath, Set: true},
			wantC: &Configuration{
				bMap: map[string]bool{},
				iMap: map[string]int{},
				sMap: map[string]string{},
				cMap: map[string]*Configuration{
					"check": {
						bMap: map[string]bool{"empty": true, "gaps": true, "integrity": false},
						iMap: map[string]int{},
						sMap: map[string]string{},
						cMap: map[string]*Configuration{},
					},
					"common": {
						bMap: map[string]bool{},
						iMap: map[string]int{},
						sMap: map[string]string{
							"albumFilter":  "^.*$",
							"artistFilter": "^.*$",
							"ext":          ".mpeg",
							"topDir":       ".",
						},
						cMap: map[string]*Configuration{},
					},
					"ls": {
						bMap: map[string]bool{
							"annotate":       true,
							"includeAlbums":  false,
							"includeArtists": false,
							"includeTracks":  true,
						},
						iMap: map[string]int{},
						sMap: map[string]string{"sort": "alpha"},
						cMap: map[string]*Configuration{},
					},
					"repair": {
						bMap: map[string]bool{"dryRun": true},
						iMap: map[string]int{},
						sMap: map[string]string{},
						cMap: map[string]*Configuration{},
					},
					"unused": {
						bMap: map[string]bool{},
						iMap: map[string]int{},
						sMap: map[string]string{"value": "1.25"},
						cMap: map[string]*Configuration{}},
				},
			},
			wantOk: true,
			WantedOutput: WantedOutput{
				WantLogOutput: fmt.Sprintf("level='warn' key='value' type='float64' value='1.25' msg='unexpected value type'\n"+
					"level='info' directory='%s' fileName='defaults.yaml' value='map[check:map[empty:true gaps:true integrity:false] common:map[albumFilter:^.*$ artistFilter:^.*$ ext:.mpeg topDir:.] ls:map[annotate:true includeAlbums:false includeArtists:false includeTracks:true], map[sort:alpha] repair:map[dryRun:true] unused:map[value:1.25]]' msg='read configuration file'\n", mp3Path),
			},
		},
		{
			name:   "APPDATA not set",
			state:  &SavedEnvVar{Name: appDataVar},
			wantC:  EmptyConfiguration(),
			wantOk: true,
			WantedOutput: WantedOutput{
				WantLogOutput: "level='info' environment variable='APPDATA' msg='not set'\n",
			},
		},
		{
			name:  "defaults.yaml is a directory",
			state: &SavedEnvVar{Name: appDataVar, Value: SecureAbsolutePathForTesting(badDir), Set: true},
			WantedOutput: WantedOutput{
				WantErrorOutput: fmt.Sprintf("The configuration file %q is a directory.\n", yamlAsDir),
				WantLogOutput:   fmt.Sprintf("level='error' directory='%s' fileName='defaults.yaml' msg='file is a directory'\n", SecureAbsolutePathForTesting(badDir2)),
			},
		},
		{
			name:   "missing yaml",
			state:  &SavedEnvVar{Name: appDataVar, Value: SecureAbsolutePathForTesting(yamlAsDir), Set: true},
			wantC:  EmptyConfiguration(),
			wantOk: true,
			WantedOutput: WantedOutput{
				WantLogOutput: fmt.Sprintf("level='info' directory='%s' fileName='defaults.yaml' msg='file does not exist'\n", SecureAbsolutePathForTesting(filepath.Join(yamlAsDir, AppName))),
			},
		},
		{
			name:  "malformed yaml",
			state: &SavedEnvVar{Name: appDataVar, Value: SecureAbsolutePathForTesting(badDir2), Set: true},
			WantedOutput: WantedOutput{
				WantErrorOutput: fmt.Sprintf(
					"The configuration file %q is not well-formed YAML: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `gibberish` into map[string]interface {}\n",
					SecureAbsolutePathForTesting(filepath.Join(gibberishDir, defaultConfigFileName))),
				WantLogOutput: fmt.Sprintf("level='warn' directory='%s' error='yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `gibberish` into map[string]interface {}' fileName='defaults.yaml' msg='cannot unmarshal yaml content'\n", SecureAbsolutePathForTesting(gibberishDir)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.state.RestoreForTesting()
			o := NewOutputDeviceForTesting()
			gotC, gotOk := ReadConfigurationFile(o)
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("%s gotC = %v, want %v", fnName, gotC, tt.wantC)
			}
			if gotOk != tt.wantOk {
				t.Errorf("%s gotOk = %v, want %v", fnName, gotOk, tt.wantOk)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}

func Test_appData(t *testing.T) {
	fnName := "appData()"
	savedState := SaveEnvVarForTesting(appDataVar)
	os.Setenv(appDataVar, SecureAbsolutePathForTesting("."))
	defer func() {
		savedState.RestoreForTesting()
	}()
	tests := []struct {
		name  string
		state *SavedEnvVar
		want  string
		want1 bool
		WantedOutput
	}{
		{
			name:         "value is set",
			state:        &SavedEnvVar{Name: appDataVar, Value: "appData!", Set: true},
			want:         "appData!",
			want1:        true,
			WantedOutput: WantedOutput{},
		},
		{
			name:  "value is not set",
			state: &SavedEnvVar{Name: appDataVar},
			WantedOutput: WantedOutput{
				WantLogOutput: "level='info' environment variable='APPDATA' msg='not set'\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.state.RestoreForTesting()
			o := NewOutputDeviceForTesting()
			got, got1 := appData(o)
			if got != tt.want {
				t.Errorf("%s got = %q, want %q", fnName, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%s got1 = %v, want %v", fnName, got1, tt.want1)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}

func TestConfiguration_StringValue(t *testing.T) {
	fnName := "Configuration.StringValue()"
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		c         *Configuration
		args      args
		wantValue string
		wantOk    bool
	}{
		{
			name:      "found",
			c:         &Configuration{sMap: map[string]string{"key": "value"}},
			args:      args{key: "key"},
			wantValue: "value",
			wantOk:    true,
		},
		{
			name: "not found",
			c:    EmptyConfiguration(),
			args: args{key: "key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.c.StringValue(tt.args.key)
			if gotValue != tt.wantValue {
				t.Errorf("%s gotValue = %q, want %q", fnName, gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("%s gotOk = %v, want %v", fnName, gotOk, tt.wantOk)
			}
		})
	}
}

func TestConfiguration_String(t *testing.T) {
	fnName := "Configuration.String()"
	tests := []struct {
		name string
		c    *Configuration
		want string
	}{
		{
			name: "empty case",
			c:    EmptyConfiguration(),
		},
		{
			name: "busy case",
			c: &Configuration{
				bMap: map[string]bool{
					"f": false,
					"t": true,
				},
				iMap: map[string]int{
					"zero": 0,
					"one":  1,
				},
				sMap: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
				cMap: map[string]*Configuration{
					"empty": EmptyConfiguration(),
				},
			},
			want: "map[f:false t:true], map[one:1 zero:0], map[bar:foo foo:bar], map[empty:]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("%s = %q, want %q", fnName, got, tt.want)
			}
		})
	}
}

func TestNewIntBounds(t *testing.T) {
	fnName := "NewIntBounds()"
	type args struct {
		v1 int
		v2 int
		v3 int
	}
	tests := []struct {
		name string
		args args
		want *IntBounds
	}{
		{
			name: "case 1",
			args: args{v1: 1, v2: 2, v3: 3},
			want: &IntBounds{minValue: 1, defaultValue: 2, maxValue: 3},
		},
		{
			name: "case 2",
			args: args{v1: 1, v2: 3, v3: 2},
			want: &IntBounds{minValue: 1, defaultValue: 2, maxValue: 3},
		},
		{
			name: "case 3",
			args: args{v1: 2, v2: 1, v3: 3},
			want: &IntBounds{minValue: 1, defaultValue: 2, maxValue: 3},
		},
		{
			name: "case 4",
			args: args{v1: 2, v2: 3, v3: 1},
			want: &IntBounds{minValue: 1, defaultValue: 2, maxValue: 3},
		},
		{
			name: "case 5",
			args: args{v1: 3, v2: 1, v3: 2},
			want: &IntBounds{minValue: 1, defaultValue: 2, maxValue: 3},
		},
		{
			name: "case 6",
			args: args{v1: 3, v2: 2, v3: 1},
			want: &IntBounds{minValue: 1, defaultValue: 2, maxValue: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIntBounds(tt.args.v1, tt.args.v2, tt.args.v3); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s = %v, want %v", fnName, got, tt.want)
			}
		})
	}
}

func TestConfiguration_IntDefault(t *testing.T) {
	fnName := "Configuration.IntDefault()"
	type args struct {
		key          string
		sortedBounds *IntBounds
	}
	tests := []struct {
		name  string
		c     *Configuration
		args  args
		wantI int
	}{
		{
			name:  "miss",
			c:     EmptyConfiguration(),
			args:  args{key: "k", sortedBounds: NewIntBounds(1, 2, 3)},
			wantI: 2,
		},
		{
			name:  "hit int value too low",
			c:     &Configuration{iMap: map[string]int{"k": -1}},
			args:  args{key: "k", sortedBounds: NewIntBounds(1, 2, 3)},
			wantI: 1,
		},
		{
			name:  "hit int value too high",
			c:     &Configuration{iMap: map[string]int{"k": 10}},
			args:  args{key: "k", sortedBounds: NewIntBounds(1, 2, 3)},
			wantI: 3,
		},
		{
			name:  "hit int value in the middle",
			c:     &Configuration{iMap: map[string]int{"k": 15}},
			args:  args{key: "k", sortedBounds: NewIntBounds(10, 20, 30)},
			wantI: 15,
		},
		{
			name:  "hit string value too low",
			c:     &Configuration{iMap: map[string]int{}, sMap: map[string]string{"k": "-1"}},
			args:  args{key: "k", sortedBounds: NewIntBounds(1, 2, 3)},
			wantI: 1,
		},
		{
			name:  "hit string value too high",
			c:     &Configuration{iMap: map[string]int{}, sMap: map[string]string{"k": "10"}},
			args:  args{key: "k", sortedBounds: NewIntBounds(1, 2, 3)},
			wantI: 3,
		},
		{
			name:  "hit string value in the middle",
			c:     &Configuration{iMap: map[string]int{}, sMap: map[string]string{"k": "15"}},
			args:  args{key: "k", sortedBounds: NewIntBounds(10, 20, 30)},
			wantI: 15,
		},
		{
			name: "hit invalid string value",
			c:     &Configuration{iMap: map[string]int{}, sMap: map[string]string{"k": "foo"}},
			args:  args{key: "k", sortedBounds: NewIntBounds(10, 20, 30)},
			wantI: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := tt.c.IntDefault(tt.args.key, tt.args.sortedBounds); gotI != tt.wantI {
				t.Errorf("%s = %d, want %d", fnName, gotI, tt.wantI)
			}
		})
	}
}

func Test_createConfiguration(t *testing.T) {
	fnName := "createConfiguration()"
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *Configuration
		WantedOutput
	}{
		{
			name: "busy!",
			args: args{
				data: map[string]interface{}{
					"boolValue":   true,
					"intValue":    1,
					"stringValue": "foo",
					"weirdValue":  1.2345,
					"mapValue":    map[string]interface{}{},
				},
			},
			want: &Configuration{
				bMap: map[string]bool{"boolValue": true},
				iMap: map[string]int{"intValue": 1},
				sMap: map[string]string{"stringValue": "foo", "weirdValue": "1.2345"},
				cMap: map[string]*Configuration{"mapValue": EmptyConfiguration()},
			},
			WantedOutput: WantedOutput{
				WantLogOutput: "level='warn' key='weirdValue' type='float64' value='1.2345' msg='unexpected value type'\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOutputDeviceForTesting()
			if got := createConfiguration(o, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s = %v, want %v", fnName, got, tt.want)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}
