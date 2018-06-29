#!/bin/bash

_bwBashNotFound() {
  local _ansiReset=$'\e[0m'
  local _ansiRed=$'\e[31m'
  local _ansiBlue=$'\e[34m'
  local _ansiUnderline=$'\e[4m'
  local _ansiUrl="${_ansiBlue}${_ansiUnderline}"
  echo "${_ansiErr}ERR: bw.bash в системе не обнаружен, см. ${_ansiUrl}https://github.com/baza-winner/bw${_ansiReset}"
  return 1
}

_bwCreateServerHttpsCrtKeyParams=( 
  '--yes/y'
)
_bwCreateServerHttpsCrtKey_description='
  Выдает список команд по cозданию корневого сертификата ${_ansiFileSpec}rootCA.pem${_ansiReset}, 
  а также файлов ${_ansiFileSpec}server.crt${_ansiReset}, ${_ansiFileSpec}server.key${_ansiReset}
'
_bwCreateServerHttpsCrtKey_yes_description='Не просто выдает список команд, а и выполняет их'
_bwCreateServerHttpsCrtKey() { eval "$_funcParams2"
  if [[ -n $yes ]]; then
    local -a OPT=( -v all -s no )
  else
    local -a OPT=( -v dry )
  fi
  # https://medium.freecodecamp.org/how-to-get-https-working-on-your-local-development-environment-in-5-minutes-7af615770eec
  _exec "${OPT[@]}" openssl genrsa -des3 -out rootCA.key 2048
  _exec "${OPT[@]}" openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 36500 -out rootCA.pem
  _exec "${OPT[@]}" openssl req -new -sha256 -nodes -out server.csr -newkey rsa:2048 -keyout server.key -config <( cat server.csr.cnf )
  _exec "${OPT[@]}" openssl x509 -req -in server.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out server.crt -days 500 -sha256 -extfile v3.ext
}

if [[ -z "$_bwFileSpec" ]]; then
  _bwBashNotFound
else
  . "$_bwFileSpec" -p -
  _bwCreateServerHttpsCrtKey "$@"
fi
