package core

import (
    "reflect"
    "testing"
)

func TestGetParams(t *testing.T) {
    regEx := `^(?P<up>\d{2,3})/(?P<down>\d{2,3})/(?P<pulse>\d{2,3})$`
    goodStr := "120/80/60"
    goodMap := map[string]string{"": "120/80/60", "up": "120", "down": "80", "pulse": "60"}

    tests := []struct {
        name          string
        regEx         string
        str           string
        wantParamsMap map[string]string
    }{
        {name: "Success", regEx: regEx, str: goodStr, wantParamsMap: goodMap},
        {name: "Bad string 1", regEx: regEx, str: "120/80", wantParamsMap: nil},
        {name: "Bad string 2", regEx: regEx, str: "120.80.60", wantParamsMap: nil},
        {name: "Bad string 3", regEx: regEx, str: "120 80 60", wantParamsMap: nil},
        {name: "Bad string 4", regEx: regEx, str: "Test string 120/80/60", wantParamsMap: nil},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if gotParamsMap := GetParams(tt.regEx, tt.str); !reflect.DeepEqual(gotParamsMap, tt.wantParamsMap) {
                t.Errorf("GetParams() = %v, want %v", gotParamsMap, tt.wantParamsMap)
            }
        })
    }
}
