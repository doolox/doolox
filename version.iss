; ============================================================================
;  Doolox version — THE single source of truth.
;
;  Bump the version here and only here; all build/packaging tooling reads it:
;    * installer.iss   (Inno Setup, via #include)
;    * build-deb.sh    (extracts it with sed, injects it into the Debian control)
;    * the Debian package (Version field, filled in by build-deb.sh)
; ============================================================================
#define MyAppVersion "1.1.0"
