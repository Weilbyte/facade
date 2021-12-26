package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	peparser "github.com/saferwall/pe"
)

const LINKER_TEMPLATE = "#pragma comment(linker,\"/export:FACADE_EXPORTNAME=FACADE_DLLNAME_o.FACADE_EXPORTNAME,@FACADE_EXPORTORDINAL\")"

// File generators

func GenerateProject(pe *peparser.File, path string, embed bool, outDir string) {
	var info GenInfo = GenInfo{
		path:      path,
		dllName:   strings.Split(filepath.Base(path), ".")[0],
		exports:   pe.Export,
		embedUUID: "",
	}

	if embed {
		info.embedUUID = getUUID()
	}

	_ = os.MkdirAll(outDir, 0755)
	os.WriteFile(filepath.Join(outDir, "CMakeLists.txt"), []byte(generateCMake(info)), 0644)
	os.WriteFile(filepath.Join(outDir, "main.cpp"), []byte(generateSource(info)), 0644)
}

func generateCMake(info GenInfo) string {
	return fmt.Sprintf(
		`cmake_minimum_required(VERSION 3.13)
project(%s_proxy)

add_library(%s SHARED main.cpp)`, info.dllName, info.dllName)
}

func generateSource(info GenInfo) string {
	var result string
	result += generateLinkerDirectives(info)
	result += "#include <windows.h>\n"
	if info.embedUUID != "" {
		result += "#include <stdio.h>\n"
	}

	result +=
		fmt.Sprintf(`
void Payload() {
	// TODO: Implement payload here
}
%s

BOOL WINAPI DllMain(HINSTANCE hinstDLL, DWORD fdwReason, LPVOID lpReserved) {
	switch (fdwReason) {
		case DLL_PROCESS_ATTACH:
			%sPayload();
			break;
		case DLL_PROCESS_DETACH:
			break;
	}
	
	return TRUE;
}
	`, generateEmbed(info), generateEmbedAttach(info))
	return result
}

// Sub-generators

func generateLinkerDirectives(info GenInfo) string {
	var result string
	var template string = LINKER_TEMPLATE

	if info.embedUUID != "" {
		template = strings.Replace(template, "FACADE_DLLNAME_o", fmt.Sprintf("\\\"C:\\\\Windows\\\\Temp\\\\%s-FACADE_DLLNAME.dll\\\"", info.embedUUID), 1)
	}

	for _, function := range info.exports.Functions {
		toAppend := strings.ReplaceAll(template, "FACADE_EXPORTNAME", function.Name)
		toAppend = strings.ReplaceAll(toAppend, "FACADE_DLLNAME", info.dllName)
		toAppend = strings.ReplaceAll(toAppend, "FACADE_EXPORTORDINAL", fmt.Sprintf("%d", function.Ordinal))
		result += fmt.Sprintf("%s\n", toAppend)
	}

	return result
}

func generateEmbed(info GenInfo) string {
	if info.embedUUID != "" {
		var result string = "\nunsigned char dll[] = {\n"
		var first bool = true
		fileBytes, _ := ioutil.ReadFile(info.path)
		for _, b := range fileBytes {
			if first {
				result += fmt.Sprintf("0x%02X", b)
				first = false
				continue
			}
			result += fmt.Sprintf(", 0x%02X", b)
		}
		result += fmt.Sprintf("};\nunsigned int dll_len = %d;", len(fileBytes))
		return result
	}
	return ""
}

func generateEmbedAttach(info GenInfo) string {
	if info.embedUUID != "" {
		return fmt.Sprintf(`FILE* file;
			fopen_s(&file, "C:\\Windows\\Temp\\%s-%s.dll", "wb");
			fwrite(dll, 1, dll_len, file);
			fclose(file);
		`, info.embedUUID, info.dllName)
	}
	return ""
}
