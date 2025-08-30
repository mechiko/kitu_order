package config

// если нужно по умолчанию имя,
// используется в Config как имя файла конфиг и в repo//DbSelf как имя БД
const Name = "kituorder"

var Mode = "development"

// This should preferably be set at build time via build scripts
// Set during build (adjust module path): go build -ldflags "-X 'kitu/config.ExeVersion=v1.0.0'"
const ExeVersion string = "0.0.1"

var DbVersion = "202504251545" // YYYYmmDDHHmm

var FsrarId = ""
