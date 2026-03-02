package core

import "regexp"

func GetParams(regEx, str string) (paramsMap map[string]string) {
    var compRegEx = regexp.MustCompile(regEx)
    match := compRegEx.FindStringSubmatch(str)

    if len(match) < 1 {
        return nil
    }

    paramsMap = make(map[string]string)
    names := compRegEx.SubexpNames()

    for i, p := range match {
        paramsMap[names[i]] = p
    }

    return paramsMap
}
