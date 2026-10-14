; ============================================================================
;  Doolox Windows installer
;  ----------------------------------------------------------------------------
;  Requires Inno Setup 6+ (https://jrsoftware.org/isinfo.php)
;
;  Build the Windows binary FIRST (from PowerShell in the project root):
;      .\build.ps1
;  which produces bin\doolox-windows-amd64.exe. Then compile this script
;  (ISCC.exe installer.iss) -> dist\doolox-setup-1.0.0.exe.
;
;  Layout after install:
;      C:\Program Files\Doolox\
;          doolox.exe
;          assets\...                         (tray icon + assets)
;      %USERPROFILE%\Websites\
;          config.json
;          doolox.db                         (pre-seeded example pages)
;          content\...                       (templates + theme assets)
;          docs\                             (empty; static site is generated here)
;
;  IMPORTANT: Doolox reads config.json, doolox.db, content\ and docs\ relative
;  to its current working directory. The Start Menu shortcut therefore sets its
;  "Start in" (WorkingDir) to %USERPROFILE%\Websites so the app works out of the
;  box without any further setup.
;
;  Finding your data: open File Explorer, click the address bar, type
;      %USERPROFILE%\Websites
;  and press Enter. That "Websites" is a real folder in your home directory
;  (C:\Users\<you>\Websites). The Start Menu / desktop "Doolox" entries are only
;  shortcuts to the doolox.exe program, not to your data. A "Doolox Websites"
;  shortcut is added to the Start Menu to open the data folder directly.
; ============================================================================

#define MyAppName "Doolox"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Doolox"
#define MyAppExeName "doolox.exe"

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
; Writable user data -> %USERPROFILE%\Websites.
; "onlyifdoesntexist" + "uninsneveruninstall" so reinstalls/upgrades and
; uninstalls never clobber or delete the user's site content and database.
; ---------------------------------------------------------------------------
Source: "content\*"; DestDir: "{%USERPROFILE}\Websites\content"; Flags: ignoreversion recursesubdirs createallsubdirs onlyifdoesntexist uninsneveruninstall
Source: "config.json"; DestDir: "{%USERPROFILE}\Websites"; Flags: ignoreversion onlyifdoesntexist uninsneveruninstall
Source: "doolox.db";  DestDir: "{%USERPROFILE}\Websites"; Flags: ignoreversion onlyifdoesntexist uninsneveruninstall

[Dirs]
; Empty docs folder: Doolox writes the generated static site here at runtime.
Name: "{%USERPROFILE}\Websites\docs"; Flags: uninsneveruninstall

[Icons]
; Start Menu entry. WorkingDir is what makes the app "just work" after install.
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{%USERPROFILE}\Websites"; IconFilename: "{app}\assets\trayicon.ico"; Comment: "Launch Doolox CMS"

; Optional desktop shortcut (same working directory).
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{%USERPROFILE}\Websites"; IconFilename: "{app}\assets\trayicon.ico"; Comment: "Launch Doolox CMS"; Tasks: desktopicon

; Shortcut that opens the Websites folder (your site data) in Explorer.
Name: "{group}\Doolox Websites"; Filename: "{%USERPROFILE}\Websites"; Comment: "Open the folder where your website content is stored"

[Run]
; Optionally launch Doolox when the installer finishes.
; "nowait" is required: it is a system-tray app that keeps running until quit.
Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{%USERPROFILE}\Websites"; Description: "Launch {#MyAppName} now"; Flags: nowait postinstall skipifsilent

; Open the Websites folder so it is easy to find after installing.
Filename: "{%USERPROFILE}\Websites"; Description: "Open the Websites folder"; Flags: postinstall shellexec skipifsilent
