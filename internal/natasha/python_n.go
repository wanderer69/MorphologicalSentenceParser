package natasha

import (
	"errors"
	"fmt"
)

/*
#cgo windows CFLAGS: -IC:/DevelopmentSupport/Python38/include -Wstrict-prototypes -Wunused-variable
#cgo windows LDFLAGS: -LC:/DevelopmentSupport/Python38/libs -lpython38
#cgo linux CFLAGS: -I/usr/include/python3.8 -Wstrict-prototypes -Wunused-variable
#cgo linux LDFLAGS: -L/usr/lib/python3.8/config-3.8-x86_64-linux-gnu -lpython3.8
#include "Python.h"
static void AddModuleLocalPath(void) {
    PyObject *sys_path = PySys_GetObject("path");
    PyList_Append(sys_path, PyUnicode_FromString("./Py"));
}
static PyObject *pName = NULL, *pModule = NULL;
static PyObject *pDict = NULL, *pObject_ = NULL, *pVal = NULL;
static PyObject* sys = NULL;
static PyObject* sys_path = NULL;
static PyObject* folder_path = NULL;
// Загрузка интерпритатора python и модуля <module_name>.py в него.
static PyObject * pythonInit(char * path) {
    // Инициализировать интерпретатор Python
    do {
        // Загрузка модуля sys
        sys = PyImport_ImportModule("sys");
        sys_path = PyObject_GetAttrString(sys, "path");
        // Путь до наших исходников python
        folder_path = PyUnicode_FromString((const char*) path);
        PyList_Append(sys_path, folder_path);

        return folder_path;
    } while (0);

    // Печать ошибки
    PyErr_Print();
    return NULL;
}
static PyObject * pythonModuleInit(char * module_name) {
    PyObject *type, *value, *traceback;
    do {
        // Создание Unicode объекта из UTF-8 строки
        pName = PyUnicode_FromString(module_name);
        if (!pName) {
            break;
        }

        // Загрузить модуль module_name
        pModule = PyImport_Import(pName);
        if (!pModule) {
            break;
        }

        // Словарь объектов содержащихся в модуле
        pDict = PyModule_GetDict(pModule);
        if (!pDict) {
            break;
        }

        return pDict;
    } while (0);

    PyErr_Fetch(&type, &value, &traceback);
    if ((type == NULL) && ((value == NULL) && (traceback == NULL))) {
        return NULL;
    }

    // Печать ошибки
    PyErr_Print();
    return value;
}
// Освобождение ресурсов интерпретатора python
static void pythonClear() {
    // Вернуть ресурсы системе
    Py_XDECREF(pDict);
    Py_XDECREF(pModule);
    Py_XDECREF(pName);

    Py_XDECREF(folder_path);
    Py_XDECREF(sys_path);
    Py_XDECREF(sys);

    // Выгрузка интерпритатора Python
    Py_Finalize();
}
// Передача строки в качестве аргумента и получение строки назад
static char * pythonFuncGetStr(char * func_name, char *val) {
    char *ret = NULL;

    // Загрузка объекта get_value из func.py
    pObject_ = PyDict_GetItemString(pDict, (const char *) func_name); // "get_value"
    if (!pObject_) {
        return ret;
    }

    do {
        // Проверка pObject_ на годность.
        if (!PyCallable_Check(pObject_)) {
            break;
        }

        pVal = PyObject_CallFunction(pObject_, (char *) "(s)", val);
        if (pVal != NULL) {
            PyObject* pResultRepr = PyObject_Repr(pVal);

            // Если полученную строку не скопировать, то после очистки ресурсов python её не будет.
            // Для начала pResultRepr нужно привести к массиву байтов.
            ret = strdup(PyBytes_AS_STRING(PyUnicode_AsEncodedString(pResultRepr, "utf-8", "ERROR")));

            Py_XDECREF(pResultRepr);
            Py_XDECREF(pVal);
        } else {
            PyErr_Print();
        }
    } while (0);

    return ret;
}
// Получение значения переменной содержащей значение типа int
static int pythonFuncGetVal(char *val) {
    int ret = 0;

    // Получить объект с именем val
    pVal = PyDict_GetItemString(pDict, (const char *) val);
    if (!pVal) {
        return ret;
    }

    // Проверка переменной на long
    if (PyLong_Check(pVal)) {
        ret = _PyLong_AsInt(pVal);
    } else {
        PyErr_Print();
    }

    return ret;
}
*/
import "C"

type PyCmd struct {
	Chain chan int
}

func Py_IsInitialized() bool {
	return C.Py_IsInitialized() != 0
}

func (pc *PyCmd) Py_cmd_Init(path string) error {
	pc.Chain = make(chan int)
	C.Py_Initialize()
	if !Py_IsInitialized() {
		return errors.New("Error initializing the python interpreter")
	}
	if len(path) == 0 {
		path = "./Py"
	}
	pathToModule := C.CString(path)
	C.pythonInit(pathToModule)
	C.AddModuleLocalPath()
	return nil
}

func (pc *PyCmd) Py_cmd_Import(moduleName string) {
	cModuleName := C.CString(moduleName)
	C.pythonModuleInit(cModuleName)
}

func (pc *PyCmd) Py_cmd_Wait() {
	// надо дождаться окончания
	for {
		data := <-pc.Chain
		if data == 1 {
			break
		}
	}
}

func (pc *PyCmd) Py_cmd_Close() {
	C.pythonClear()
}

func (pc *PyCmd) Py_cmd_Call(funcName string, args []interface{}) string {
	cFuncName := C.CString(funcName)
	arg := ""
	for i := range args {
		switch v := args[i].(type) {
		case string:
			arg = arg + string(v)
		case int32, int64:
			n := v.(int)
			arg = arg + fmt.Sprintf("%v", n)
		default:
			fmt.Println("unknown")
		}
	}
	arg_c := C.CString(arg)

	cResult := C.pythonFuncGetStr(cFuncName, arg_c)
	result := C.GoString(cResult)
	return result
}

/*
func pythonRepr(o *python3.PyObject) (string, error) {
	if o == nil {
		return "", fmt.Errorf("object is nil")
	}

	s := o.Repr()
	if s == nil {
		python3.PyErr_Clear()
		return "", fmt.Errorf("failed to call Repr object method")
	}
	defer s.DecRef()

	return python3.PyUnicode_AsUTF8(s), nil
}
*/
