; ============================================================================
;  Doolox Windows installer
;  ----------------------------------------------------------------------------
;  Requires Inno Setup 6+ (https://jrsoftware.org/isinfo.php)
;
;  Build the Windows binary FIRST (from PowerShell in the project root):
;      .\build.ps1
;  which produces bin\doolox-windows-amd64.exe. Then compile this script
;  (ISCC.exe installer.iss) -> dist\doolox-setup-<version>.exe.
;
;  Layout after install:
;      C:\Program Files\Doolox\
;          doolox.exe
;          assets\...                         (tray icon + assets)
;          admin\theme\...                    (admin web interface)
;      %USERPROFILE%\doolox\
;          config.json
;          websites\
;              content\...                   (templates + theme assets)
;              docs\                         (empty; static site is generated here)
;          templates\                        (empty; reserved for templates)
;          doolox.log                        (created by the app at runtime)
;
;  IMPORTANT: Doolox resolves its data directory (%USERPROFILE%\doolox) on its
;  own, so it works out of the box regardless of the shortcut's working folder.
;
;  Finding your data: open File Explorer, click the address bar, type
;      %USERPROFILE%\doolox
;  and press Enter. The Start Menu / desktop "Doolox" entries are only
;  shortcuts to the doolox.exe program, not to your data. A "Doolox data"
;  shortcut is added to the Start Menu to open the data folder directly.
; ============================================================================

#define MyAppName "Doolox"
#define MyAppPublisher "Doolox"
#define MyAppExeName "doolox.exe"

; App version lives in version.iss (single source of truth for all targets).
#include "version.iss"

[Setup]
; Unique, stable identifier for this application. Keep it the same across
; releases so upgrades are detected correctly.
AppId={{9992E0D0-4D05-4108-A38C-277DAE62BFD0}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\Doolox
DefaultGroupName=Doolox
DisableProgramGroupPage=yes
UninstallDisplayName={#MyAppName}
UninstallDisplayIcon={app}\assets\trayicon.ico

; Installer icon (relative to this .iss file).
SetupIconFile=assets\trayicon.ico

; Install to the 64-bit Program Files folder (the binary is windows/amd64).
PrivilegesRequired=admin
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64

OutputDir=dist
OutputBaseFilename=doolox-setup-{#MyAppVersion}

Compression=lzma2
SolidCompression=yes
WizardStyle=modern

[Tasks]
; Optional desktop shortcut in addition to the Start Menu entry.
Name: "desktopicon"; Description: "Create a &desktop shortcut"; GroupDescription: "Additional shortcuts:"; Flags: unchecked

[Files]
; ---------------------------------------------------------------------------
; The application binary (renamed to a friendly name).
; ---------------------------------------------------------------------------
Source: "bin\doolox-windows-amd64.exe"; DestDir: "{app}"; DestName: "{#MyAppExeName}"; Flags: ignoreversion

; ---------------------------------------------------------------------------
; Project assets (tray icon etc.), copied alongside the binary.
; ---------------------------------------------------------------------------
Source: "assets\*"; DestDir: "{app}\assets"; Flags: ignoreversion recursesubdirs createallsubdirs

; ---------------------------------------------------------------------------
; Admin web interface (theme templates + static assets served on :5555).
; The admin server resolves admin\theme\static relative to the executable.
; ---------------------------------------------------------------------------
Source: "admin\theme\*"; DestDir: "{app}\admin\theme"; Flags: ignoreversion recursesubdirs createallsubdirs

; ---------------------------------------------------------------------------
; Writable user data -> %USERPROFILE%\doolox.
; "onlyifdoesntexist" + "uninsneveruninstall" so reinstalls/upgrades and
; uninstalls never clobber or delete the user's site content and config.
; ---------------------------------------------------------------------------
Source: "defaults\doolox\content\*"; DestDir: "{%USERPROFILE}\doolox\websites\doolox\content"; Flags: ignoreversion recursesubdirs createallsubdirs onlyifdoesntexist uninsneveruninstall
Source: "defaults\config.json"; DestDir: "{%USERPROFILE}\doolox"; Flags: ignoreversion onlyifdoesntexist uninsneveruninstall

[Dirs]
; Empty docs folder (generated static site) and the reserved templates folder.
Name: "{%USERPROFILE}\doolox\websites\doolox\docs"; Flags: uninsneveruninstall
Name: "{%USERPROFILE}\doolox\templates"; Flags: uninsneveruninstall

[Icons]
; Start Menu entry.
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{%USERPROFILE}\doolox"; IconFilename: "{app}\assets\trayicon.ico"; Comment: "Launch Doolox CMS"

; Optional desktop shortcut.
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{%USERPROFILE}\doolox"; IconFilename: "{app}\assets\trayicon.ico"; Comment: "Launch Doolox CMS"; Tasks: desktopicon

; Shortcut that opens the data folder in Explorer.
Name: "{group}\Doolox data"; Filename: "{%USERPROFILE}\doolox"; Comment: "Open the folder where your website data is stored"

[Run]
; Optionally launch Doolox when the installer finishes.
; "nowait" is required: it is a system-tray app that keeps running until quit.
Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{%USERPROFILE}\doolox"; Description: "Launch {#MyAppName} now"; Flags: nowait postinstall skipifsilent

; Open the data folder so it is easy to find after installing.
Filename: "{%USERPROFILE}\doolox"; Description: "Open the Doolox data folder"; Flags: postinstall shellexec skipifsilent
