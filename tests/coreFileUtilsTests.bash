
# =============================================================================

_resetBash

# =============================================================================

_downloadTestsCondition='[[ -n $_isBwDevelop || -n $BW_TEST_ALL ]]'
_downloadTests=(
  # '
  #   --return "2"
  #   "--stdout 
  #     ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_download${_ansiReset} [${_ansiOutline}Опции${_ansiReset}] ${_ansiOutline}url${_ansiReset} [${_ansiOutline}fileSpec${_ansiReset}]
  #     ${_ansiHeader}Описание:${_ansiReset}
  #       Загружает содержимое ${_ansiOutline}url${_ansiReset}'\''а в файл ${_ansiOutline}fileSpec${_ansiReset}
  #       Если загрузка была прервана, то продолжает загрузку
  #     ${_ansiOutline}Опции${_ansiReset}
  #       ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset} Выводит справку
  #       ${_ansiCmd}--silent${_ansiReset} или ${_ansiCmd}-s${_ansiReset} Молчаливый режим
  #       ${_ansiCmd}--check ${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-c ${_ansiOutline}значение${_ansiReset}
  #         ${_ansiOutline}значение${_ansiReset} - один из вариантов: ${_ansiSecondaryLiteral}etag time force${_ansiReset}
  #         ${_ansiPrimaryLiteral}etag${_ansiReset}
  #           Загружает файл, если файл ранее не был загружен,
  #           или если на сервере изменялся ETag файла с момента последней загрузки
  #         ${_ansiPrimaryLiteral}time${_ansiReset}
  #           Загружает файл, если файл ранее был не загружен,
  #           или если на сервере он изменился с момента последней загрузки
  #         ${_ansiPrimaryLiteral}force${_ansiReset}
  #           Загружает файл с самого начала независимо ни от чего
  #           Игнорирует опцию ${_ansiCmd}--tolerant${_ansiReset}
  #       ${_ansiCmd}--tolerant${_ansiReset} или ${_ansiCmd}-t${_ansiReset}
  #         Игнорирует ошибки, связанные с качеством интернета (его может просто не быть),
  #         если задана опция ${_ansiCmd}--check${_ansiReset} и файл ${_ansiOutline}fileSpec${_ansiReset} уже существует (был ранее скачан)
  #       ${_ansiCmd}--return-code-if-actually-updated=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-r ${_ansiOutline}значение${_ansiReset}
  #         ${_ansiOutline}значение${_ansiReset} - неотрицательное целое число
  #         Устанавливает код возврата на случай,
  #         если файл ${_ansiOutline}fileSpec${_ansiReset} был действительно обновлен
  #         ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}0${_ansiReset}
  #     ${_ansiOutline}url${_ansiReset} URL для загрузки
  #     ${_ansiOutline}fileSpec${_ansiReset} Путь к файлу, куда поместить результат загрузки
  #       ${_ansiOutline}fileSpec${_ansiReset} по умолчанию - значение выражения ${_ansiOutline}$ (basename $ url)${_ansiReset}
  #   "
  #   --stdoutParaWithIndent "0"
  #   --stdoutEtaPreProcess "perl -pe s/\\\$\s/\\\$/g"
  #   "_download -?"
  # '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    "--inTmpDir"
    "_exist -n bw.bash && _download localhost:$_bwdevDockerHttp/bw.bash bw.bash && cmp bw.bash '$_stqq'$_bwDir/bw.bash'$_stqq'"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiCmd}rm -f bw.bash.download.header${_ansiReset}
      ${_ansiCmd}curl -o bw.bash.download -L localhost:$_bwdevDockerHttp/bw.bash -s --dump-header bw.bash.download.header${_ansiReset}
      ${_ansiCmd}mv bw.bash.download bw.bash${_ansiReset}
      ${_ansiCmd}mv bw.bash.download.header bw.bash.header${_ansiReset}
    "
    "--inTmpDir"
    "_exist -n bw.bash && _download -v dry localhost:$_bwdevDockerHttp/bw.bash bw.bash && _exist -n bw.bash"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    "--inTmpDir"
    "_exist -n bw.bash && _download -v none localhost:$_bwdevDockerHttp/bw.bash bw.bash && cmp bw.bash '$_stqq'$_bwDir/bw.bash'$_stqq'"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    "--inTmpDir"
    "_exist -n bw.bash && _download -v err localhost:$_bwdevDockerHttp/bw.bash bw.bash && cmp bw.bash '$_stqq'$_bwDir/bw.bash'$_stqq'"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}rm -f bw.bash.download.header${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}curl -o bw.bash.download -L localhost:$_bwdevDockerHttp/bw.bash -s --dump-header bw.bash.download.header${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mv bw.bash.download bw.bash${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mv bw.bash.download.header bw.bash.header${_ansiReset}
    "
    "--inTmpDir"
    "_exist -n bw.bash && _download -v ok localhost:$_bwdevDockerHttp/bw.bash bw.bash && cmp bw.bash '$_stqq'$_bwDir/bw.bash'$_stqq'"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}rm -f bw.bash.download.header${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}curl -o bw.bash.download -L localhost:$_bwdevDockerHttp/bw.bash -s --dump-header bw.bash.download.header${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mv bw.bash.download bw.bash${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mv bw.bash.download.header bw.bash.header${_ansiReset}
    "
    "--inTmpDir"
    "_exist -n bw.bash && _download -v allBrief localhost:$_bwdevDockerHttp/bw.bash bw.bash && cmp bw.bash '$_stqq'$_bwDir/bw.bash'$_stqq'"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}rm -f bw.bash.download.header${_ansiReset}
      ${_ansiCmd}curl -o bw.bash.download -L localhost:$_bwdevDockerHttp/bw.bash -s --dump-header bw.bash.download.header${_ansiReset} . . .
      ${_ansiOK}OK: ${_ansiCmd}curl -o bw.bash.download -L localhost:$_bwdevDockerHttp/bw.bash -s --dump-header bw.bash.download.header${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mv bw.bash.download bw.bash${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mv bw.bash.download.header bw.bash.header${_ansiReset}
    "
    "--inTmpDir"
    "_exist -n bw.bash && _download -v all localhost:$_bwdevDockerHttp/bw.bash bw.bash && cmp bw.bash '$_stqq'$_bwDir/bw.bash'$_stqq'"
  '
  '
    --before "_inDir \"$_bwDir/docker\" _dockerCompose up -d"
    "--inTmpDir"
    "_download localhost:$_bwdevDockerHttp/bw.bash bw.bash && _download -c etag localhost:$_bwdevDockerHttp/bw.bash bw.bash"
  '
)

_inDirTests=(
  '
    --stdout "$_bwTmpDir/inner/inner2"
    "--inTmpDir"
    "_inDir inner/inner2 pwd"
  '
  '
    "--inTmpDir"
    "_inDir -s yes inner/inner2 pwd"
  '
  '
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
      ${_ansiCmd}pushd inner/inner2 >/dev/null 2>&1${_ansiReset}
      ${_ansiCmd}pwd${_ansiReset}
      ${_ansiCmd}popd >/dev/null 2>&1${_ansiReset}
    "
    "--inTmpDir"
    "_inDir -v dry inner/inner2 pwd"
  '
  '
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
      ${_ansiCmd}pushd inner/inner2 >/dev/null 2>&1${_ansiReset}
      ${_ansiCmd}pwd >/dev/null 2>&1${_ansiReset}
      ${_ansiCmd}popd >/dev/null 2>&1${_ansiReset}
    "
    "--inTmpDir"
    "_inDir -v dry -s yes inner/inner2 pwd"
  '
  '
    --stdout "$_bwTmpDir/inner/inner2"
    "--inTmpDir"
    "_inDir -v none inner/inner2 pwd"
  '
  '
    "--inTmpDir"
    "_inDir -v none -s yes inner/inner2 pwd"
  '
  '
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}pushd inner/inner2 >/dev/null 2>&1${_ansiReset}
      $_bwTmpDir/inner/inner2
      ${_ansiOK}OK: ${_ansiCmd}pwd${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}popd >/dev/null 2>&1${_ansiReset}
    "
    "--inTmpDir"
    "_inDir -v ok inner/inner2 pwd"
  '
  '
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}pushd inner/inner2 >/dev/null 2>&1${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}pwd >/dev/null 2>&1${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}popd >/dev/null 2>&1${_ansiReset}
    "
    "--inTmpDir"
    "_inDir -v ok -s yes inner/inner2 pwd"
  '
  '
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}pushd inner/inner2 >/dev/null 2>&1${_ansiReset}
      ${_ansiCmd}pwd${_ansiReset} . . .
      $_bwTmpDir/inner/inner2
      ${_ansiOK}OK: ${_ansiCmd}pwd${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}popd >/dev/null 2>&1${_ansiReset}
    "
    "--inTmpDir"
    "_inDir -v all inner/inner2 pwd"
  '
  '
    --stdoutParaWithIndent "0"
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}pushd inner/inner2 >/dev/null 2>&1${_ansiReset}
      ${_ansiCmd}pwd >/dev/null 2>&1${_ansiReset} . . .
      ${_ansiOK}OK: ${_ansiCmd}pwd >/dev/null 2>&1${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}popd >/dev/null 2>&1${_ansiReset}
    "
    "--inTmpDir"
    "_inDir -v all -s yes inner/inner2 pwd"
  '
  # ===
  '
    --stdout "some.file"
    "--inTmpDir"
    --before "_mkDir inner"
    --before "touch inner/some.file"
    "_inDir inner ls . && _exist inner/some.file"
  '
  '
    "--inTmpDir"
    --before "_mkDir inner"
    --before "touch inner/some.file"
    "_inDir -t inner ls . && _exist -n inner/some.file && _exist -n -d inner"
  '
  '
    "--inTmpDir"
    --before "_mkDir inner"
    --before "touch inner/some.file"
    "_inDir -t -n -v none inner _inDirTestHelper || _exist -n inner/some.file && _exist inner/another.file"
  '
  # ===
)
_inDirTestHelper() {
  touch another.file && false
}

_rmTests=(
  '
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm some.file && _exist -n some.file"
  '
  '
    --stdout "${_ansiCmd}rm -f some.file${_ansiReset}"
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v dry some.file && _exist some.file"
  '
  '
    --stdout "${_ansiCmd}rm -f some.file >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v dry -s yes some.file && _exist some.file"
  '
  '
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v none some.file && _exist -n some.file"
  '
  '
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v err some.file && _exist -n some.file"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f some.file${_ansiReset}"
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v ok some.file && _exist -n some.file"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f some.file${_ansiReset}"
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v allBrief some.file && _exist -n some.file"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f some.file${_ansiReset}"
    "--inTmpDir"
    --before "touch some.file"
    "_exist some.file && _rm -v all some.file && _exist -n some.file"
  '
  # ===
  '
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -d inner && _exist -n -d inner"
  '
  '
    --stdout "${_ansiCmd}rm -rfd inner/${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v dry -d inner && _exist -d inner"
  '
  '
    --stdout "${_ansiCmd}rm -rfd inner/ >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v dry -s yes -d inner && _exist -d inner"
  '
  '
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v none -d inner && _exist -n -d inner"
  '
  '
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v err -d inner && _exist -n -d inner"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v ok -d inner && _exist -n -d inner"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v allBrief -d inner && _exist -n -d inner"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner"
    "_exist -d inner && _rm -v all -d inner && _exist -n -d inner"
  '
  # ===
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f symlinkTo.file${_ansiReset}"
    "--inTmpDir"
    --before "touch some.file"
    --before "ln some.file symlinkTo.file"
    "_exist some.file symlinkTo.file && _rm -v all symlinkTo.file && _exist -n symlinkTo.file && _exist some.file"
  '
  # ===
  '
    "_exist -n nonExistent && _rm nonExistent"
  '
  '
    --stdout "${_ansiCmd}rm -f nonExistent${_ansiReset}"
    "_exist -n nonExistent && _rm -v dry nonExistent"
  '
  '
    --stdout "${_ansiCmd}rm -f nonExistent >/dev/null 2>&1${_ansiReset}"
    "_exist -n nonExistent && _rm -v dry -s yes nonExistent"
  '
  '
    "_exist -n nonExistent && _rm -v none nonExistent"
  '
  '
    "_exist -n nonExistent && _rm -v err nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f nonExistent${_ansiReset}"
    "_exist -n nonExistent && _rm -v ok nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f nonExistent${_ansiReset}"
    "_exist -n nonExistent && _rm -v allBrief nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -f nonExistent${_ansiReset}"
    "_exist -n nonExistent && _rm -v all nonExistent"
  '
  # ===
  '
    "_exist -n -d nonExistent && _rm -d nonExistent"
  '
  '
    --stdout "${_ansiCmd}rm -rfd nonExistent/${_ansiReset}"
    "_exist -n -d nonExistent && _rm -v dry -d nonExistent"
  '
  '
    --stdout "${_ansiCmd}rm -rfd nonExistent/ >/dev/null 2>&1${_ansiReset}"
    "_exist -n -d nonExistent && _rm -v dry -s yes -d nonExistent"
  '
  '
    "_exist -n -d nonExistent && _rm -v none -d nonExistent"
  '
  '
    "_exist -n -d nonExistent && _rm -v err -d nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -rfd nonExistent/${_ansiReset}"
    "_exist -n -d nonExistent && _rm -v ok -d nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -rfd nonExistent/${_ansiReset}"
    "_exist -n -d nonExistent && _rm -v allBrief -d nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}rm -rfd nonExistent/${_ansiReset}"
    "_exist -n -d nonExistent && _rm -v all -d nonExistent"
  '
  # ===

  # '
  #   --return "1"
  #   "--stderr ${_ansiErr}ERR: ${_ansiCmd}rm -rfd some.file/ >/dev/null 2>&1${_ansiReset}"
  #   "--inTmpDir"
  #   --before "touch some.file"
  #   "_rm -s yes -d some.file"
  # '
  # UBUNTU:
  # ERR: _rm:33: _rm -s yes -d some.file
  # return is expected to be 1, but got
  # stderr (diff /tmp/testsSupport.bash.stderr /tmp/testsSupport.bash.stderr.eta):
  # 0a1
  # > ERR: rm -rfd some.file/ >/dev/null 2>&1

  # '
  #   --return "1"
  #   "--inTmpDir"
  #   --before "touch some.file"
  #   "_rm -v ok -s yes -d some.file"
  # '
  # ERR: _rm:34: _rm -v ok -s yes -d some.file
  # return is expected to be 1, but got
  # stdout (diff /tmp/testsSupport.bash.stdout /tmp/testsSupport.bash.stdout.eta):
  # 1d0
  # < OK: rm -rfd some.file/ >/dev/null 2>&1

  # ===
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}rm -f inner >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner"
    "_rm -s yes inner"
  '
  '
    --return "1"
    "--inTmpDir"
    --before "_mkDir inner"
    "_rm -v ok -s yes inner"
  '
  # ===
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}rm -f inner/inner2/readOnly.file >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner/inner2"
    --before "touch inner/inner2/readOnly.file"
    --before "cd inner/inner2"
    --before "chmod 477 ."
    --before "cd ../.."
    "_rm -s yes inner/inner2/readOnly.file"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}rm -rfd inner/inner2/ >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    --before "_mkDir inner/inner2"
    --before "touch inner/inner2/readOnly.file"
    --before "cd inner/inner2"
    --before "chmod 477 ."
    --before "cd ../.."
    "_exist -d inner/inner2 && _rm -s yes -d inner/inner2"
  '
)

_mkDirTests=(
  '
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}"
    "--inTmpDir"
    "_mkDir -v all inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}"
    "--inTmpDir"
    "_mkDir -v allBrief inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}"
    "--inTmpDir"
    "_mkDir -v ok inner/inner2 && _exist -d inner/inner2"
  '
  '
    "--inTmpDir"
    "_mkDir -v none inner/inner2 && _exist -d inner/inner2"
  '
  '
    "--inTmpDir"
    "_mkDir -v err inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "${_ansiCmd}mkdir -p inner/inner2${_ansiReset}"
    "--inTmpDir"
    "_mkDir -v dry inner/inner2 && _exist -d -n inner/inner2"
  '
  '
    --stdout "${_ansiCmd}mkdir -p inner/inner2 >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mkDir -v dry -s yes inner/inner2 && _exist -d -n inner/inner2"
  '

  '
    "--inTmpDir"
    "_mkDir -t inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir -v all -t inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir -v allBrief -t inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir -v ok -t inner/inner2 && _exist -d inner/inner2"
  '
  '
    "--inTmpDir"
    "_mkDir -v none -t inner/inner2 && _exist -d inner/inner2"
  '
  '
    "--inTmpDir"
    "_mkDir -v err -t inner/inner2 && _exist -d inner/inner2"
  '
  '
    --stdout "
      ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir -v dry -t inner/inner2 && _exist -d -n inner/inner2"
  '
  '
    --stdout "
      ${_ansiCmd}rm -rfd inner/inner2/ >/dev/null 2>&1${_ansiReset}
      ${_ansiCmd}mkdir -p inner/inner2 >/dev/null 2>&1${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir -v dry -s yes -t inner/inner2 && _exist -d -n inner/inner2"
  '

  '
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir inner/inner2"
  '
  '
    --stdout "${_ansiCmd}mkdir -p inner/inner2${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v dry inner/inner2"
  '
  '
    --stdout "${_ansiCmd}mkdir -p inner/inner2 >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v dry -s yes inner/inner2"
  '
  '
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v none inner/inner2"
  '
  '
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v err inner/inner2"
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}inner/inner2${_ansiOK} существует${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v all inner/inner2"
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}inner/inner2${_ansiOK} существует${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v allBrief inner/inner2"
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}inner/inner2${_ansiOK} существует${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && _exist -d inner/inner2 && _mkDir -v ok inner/inner2"
  '

  '
    "--inTmpDir"
    "_mkDir inner/inner2 && touch inner/inner2/some.file && _mkDir -t inner/inner2 && _exist -n inner/inner2/some.file"
  '
  '
    "--inTmpDir"
    "_mkDir inner/inner2 && touch inner/inner2/some.file && _mkDir -t -v none inner/inner2 && _exist -n inner/inner2/some.file"
  '
  '
    "--inTmpDir"
    "_mkDir inner/inner2 && touch inner/inner2/some.file && _mkDir -t -v err inner/inner2 && _exist -n inner/inner2/some.file"
  '
  '
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}chmod 777 inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir inner/inner2 && touch inner/inner2/some.file && _mkDir -t -v all inner/inner2 && _exist -n inner/inner2/some.file"
  '
  '
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}chmod 777 inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir inner/inner2 && touch inner/inner2/some.file && _mkDir -t -v allBrief inner/inner2 && _exist -n inner/inner2/some.file"
  '
  '
    --stdout "
      ${_ansiOK}OK: ${_ansiCmd}chmod 777 inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}rm -rfd inner/inner2/${_ansiReset}
      ${_ansiOK}OK: ${_ansiCmd}mkdir -p inner/inner2${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "--inTmpDir"
    "_mkDir inner/inner2 && touch inner/inner2/some.file && _mkDir -t -v ok inner/inner2 && _exist -n inner/inner2/some.file"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}mkdir -p inner3 >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && cd inner/inner2 && chmod 477 . && _mkDir -s yes inner3"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}mkdir -p inner3 >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && cd inner/inner2 && chmod 477 . && _mkDir -s yes -v err inner3"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}mkdir -p inner3 >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && cd inner/inner2 && chmod 477 . && _mkDir -s yes -v allBrief inner3"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}mkdir -p inner3 >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mkDir inner/inner2 && cd inner/inner2 && chmod 477 . && _mkDir -s yes -v all inner3"
  '
  '
    --return "1"
    "--inTmpDir"
    "_mkDir inner/inner2 && cd inner/inner2 && chmod 477 . && _mkDir -s yes -v ok inner3"
  '
  '
    --return "1"
    "--inTmpDir"
    "_mkDir inner/inner2 && cd inner/inner2 && chmod 477 . && _mkDir -s yes -v none inner3"
  '
)

_existTests=(
  '
    "_exist \"$_bwFileSpec\""
  '
  '
    --stdout "${_ansiOK}OK: Файл ${_ansiFileSpec}$_bwFileSpec${_ansiOK} существует${_ansiReset}"
    "_exist -v all \"$_bwFileSpec\""
  '
  '
    --stdout "${_ansiOK}OK: Файл ${_ansiFileSpec}$_bwFileSpec${_ansiOK} существует${_ansiReset}"
    "_exist -v allBrief \"$_bwFileSpec\""
  '
  '
    "_exist -n nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: Файл ${_ansiFileSpec}nonExistent${_ansiOK} не существует${_ansiReset}"
    "_exist -n -v all nonExistent"
  '
  '
    --stdout "${_ansiOK}OK: Файл ${_ansiFileSpec}nonExistent${_ansiOK} не существует${_ansiReset}"
    "_exist -n -v allBrief nonExistent"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Файл ${_ansiFileSpec}$_bwFileSpec${_ansiErr} существует${_ansiReset}"
    "_exist -n \"$_bwFileSpec\""
  '
  '
    --return "1"
    "_exist -n -v none \"$_bwFileSpec\""
  '
  '
    --return "1"
    "_exist -n -v dry \"$_bwFileSpec\""
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Файл ${_ansiFileSpec}nonExistent${_ansiErr} не существует${_ansiReset}"
    "_exist nonExistent"
  '
  '
    --return "1"
    "_exist -v none nonExistent"
  '
  '
    --return "1"
    "_exist -v dry nonExistent"
  '
  '
    "_exist -d \"$_bwDir\""
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}$_bwDir${_ansiOK} существует${_ansiReset}"
    "_exist -d -v all \"$_bwDir\""
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}$_bwDir${_ansiOK} существует${_ansiReset}"
    "_exist -d -v allBrief \"$_bwDir\""
  '
  '
    "_exist -nd \"$_bwFileSpec\""
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}$_bwFileSpec${_ansiOK} не существует${_ansiReset}"
    "_exist -nd -v all \"$_bwFileSpec\""
  '
  '
    --stdout "${_ansiOK}OK: Директория ${_ansiFileSpec}$_bwFileSpec${_ansiOK} не существует${_ansiReset}"
    "_exist -nd -v allBrief \"$_bwFileSpec\""
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Директория ${_ansiFileSpec}$_bwFileSpec${_ansiErr} не существует${_ansiReset}"
    "_exist -d \"$_bwFileSpec\""
  '
  '
    --return "1"
    "_exist -d -v none \"$_bwFileSpec\""
  '
  '
    --return "1"
    "_exist -d -v dry \"$_bwFileSpec\""
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Директория ${_ansiFileSpec}$_bwDir${_ansiErr} существует${_ansiReset}"
    "_exist -nd \"$_bwDir\""
  '
  '
    --return "1"
    "_exist -nd -v none \"$_bwDir\""
  '
  '
    --return "0"
    "_exist -v none \"$_bwFileSpec\" \"$_bwDir/$_bwFileName\""
  '
  '
    --return "0"
    "_exist -v dry \"$_bwFileSpec\" \"$_bwDir/$_bwFileName\""
  '
  '
    --return "0"
    "_exist -l any -v none \"$_bwFileSpec\" nonExistent"
  '
  '
    --return "0"
    "_exist -l any -v dry \"$_bwFileSpec\" nonExistent"
  '
)

_silentTests=(
  '
    "_silent _ok a"
  '
  '
    --return "1"
    "_silent _err a"
  '
)

_mvFileTests=(
  '
    "--inTmpDir"
    "touch some.file && _mvFile some.file thing.file && _exist -n some.file && _exist thing.file"
  '
  '
    --stdout "${_ansiCmd}mv some.file thing.file${_ansiReset}"
    "--inTmpDir"
    "touch some.file && _mvFile -v dry some.file thing.file && _exist some.file && _exist -n thing.file"
  '
  '
    --stdout "${_ansiCmd}mv some.file thing.file >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "touch some.file && _mvFile -v dry -s yes some.file thing.file && _exist some.file && _exist -n thing.file"
  '
  '
    "--inTmpDir"
    "touch some.file && _mvFile -v none some.file thing.file && _exist -n some.file && _exist thing.file"
  '
  '
    "--inTmpDir"
    "touch some.file && _mvFile -v err some.file thing.file && _exist -n some.file && _exist thing.file"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}mv some.file thing.file${_ansiReset}"
    "--inTmpDir"
    "touch some.file && _mvFile -v ok some.file thing.file && _exist -n some.file && _exist thing.file"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}mv some.file thing.file${_ansiReset}"
    "--inTmpDir"
    "touch some.file && _mvFile -v allBrief some.file thing.file && _exist -n some.file && _exist thing.file"
  '
  '
    --stdout "${_ansiOK}OK: ${_ansiCmd}mv some.file thing.file${_ansiReset}"
    "--inTmpDir"
    "touch some.file && _mvFile -v all some.file thing.file && _exist -n some.file && _exist thing.file"
  '

  '
    --return "1"
    --stderr "${_ansiErr}ERR: Файл ${_ansiFileSpec}some.file${_ansiErr} не существует${_ansiReset}"
    "--inTmpDir"
    "_mvFile some.file thing.file"
  '
  '
    --stdout "${_ansiCmd}mv some.file thing.file${_ansiReset}"
    "--inTmpDir"
    "_mvFile -v dry some.file thing.file"
  '
  '
    --stdout "${_ansiCmd}mv some.file thing.file >/dev/null 2>&1${_ansiReset}"
    "--inTmpDir"
    "_mvFile -v dry -s yes some.file thing.file"
  '
  '
    --return "1"
    "--inTmpDir"
    "_mvFile -v none some.file thing.file"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Файл ${_ansiFileSpec}some.file${_ansiErr} не существует${_ansiReset}"
    "--inTmpDir"
    "_mvFile -v err some.file thing.file"
  '
  '
    --return "1"
    "--inTmpDir"
    "_mvFile -v ok some.file thing.file"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Файл ${_ansiFileSpec}some.file${_ansiErr} не существует${_ansiReset}"
    "--inTmpDir"
    "_mvFile -v allBrief some.file thing.file"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: Файл ${_ansiFileSpec}some.file${_ansiErr} не существует${_ansiReset}"
    "--inTmpDir"
    "_mvFile -v all some.file thing.file"
  '
)

# =============================================================================

_mkFileFromTemplateTests=(
  '
    --before "echo '\''\${SOME_VAR}:\${SOME_VAR2}:\${SOME_VAR3}'\'' > some.template"
    --before "echo some_var_value:some_var2_value:some_var3_value > some.eta"
    --inTmpDir
    "SOME_VAR=some_var_value SOME_VAR2=some_var2_value SOME_VAR3=some_var3_value _mkFileFromTemplate -n some && diff some some.eta"
  '
  '
    --before "echo '\''\${SOME_VAR}:\${SOME_VAR2}:\${SOME_VAR3}'\'' > some.template"
    --before "echo '\''\${SOME_VAR}:some_var2_value:\${SOME_VAR3}'\'' > some.eta"
    --inTmpDir
    "SOME_VAR=some_var_value SOME_VAR2=some_var2_value SOME_VAR3=some_var3_value _mkFileFromTemplate -n -v SOME_VAR2 some && diff some some.eta"
  '
  '
    --before "echo '\''\${SOME_VAR}:\${SOME_VAR2}:\${SOME_VAR3}'\'' > some.template"
    --before "echo '\''some_var_value:\${SOME_VAR2}:some_var3_value'\'' > some.eta"
    --inTmpDir
    "SOME_VAR=some_var_value SOME_VAR2=some_var2_value SOME_VAR3=some_var3_value _mkFileFromTemplate -n -v SOME_VAR -v SOME_VAR3 some && diff some some.eta"
  '
)

# =============================================================================
