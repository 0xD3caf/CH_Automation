#pragma once
#include <Windows.h>

#ifdef UNICODE
#define LOADLIBRARY "LoadLibraryW"
#define strcomparei lstrcmpi
#define strlenx lstrlen
#define soc 2
#else 
#define LOADLIBRARY "LoadLibraryA"
#define strcomparei strcmpi
#define strlenx strlen
#define soc 1
#endif 

class BasicInjector {
    public:
        DWORD GetPIDFromName(const LPCTSTR name); // find process ID from process name (i.e. notepad.exe -> 348)
        DWORD CountPIDFromName(const LPCTSTR name); // count all process IDs from process name
        DWORD GetPIDListFromName(const LPCTSTR name, DWORD list[]); // get all process IDs for all running instances of i.e. notepad.exe
        void SetDebugPrivileges(BOOL state); // adjust the current process' privileges to allow us to inject code externally

        HMODULE Inject(DWORD processId, const LPCTSTR szDllPath); // inject a DLL into another process
        void Eject(DWORD processId, HMODULE module); // Unload (eject) a DLL from another process
};

DWORD BasicInjector::GetPIDFromName(const LPCTSTR name) {
    DWORD pid = 0;
    HANDLE hSnapshot;
    PROCESSENTRY32 pi;
    BOOL p32res;

    hSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (hSnapshot != INVALID_HANDLE_VALUE) {
        pi.dwSize = sizeof(PROCESSENTRY32);

        // Fetch first process from the snapshot
        p32res = Process32First(hSnapshot, &pi);

        while (p32res) {
            if (strcomparei(pi.szExeFile, name) == 0) { // found it!
                pid = pi.th32ProcessID;
                break;
            }

            // Get the next process from the snapshot
            p32res = Process32Next(hSnapshot, &pi);
        }

        CloseHandle(hSnapshot);
    }

    return pid;
}
