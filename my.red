Red []
funcName: system/options/args/1
funcParamsOpt: system/options/args/2
funcParams: copy next next system/options/args

probe funcName
probe funcParamsOpt
probe funcParams

optDefs: copy []
argDefs: copy []
varNames: copy []
letter: charset [ #"_" #"a" - #"z" #"A" - #"Z" ]
digit: charset [ #"0" - #"9" ]
foreach param funcParams [
  isArg: false
  probe parse param [ 
    [
      "--" ( isArg: true optDefs: append optDefs param ) | 
      ( argDefs: append argDefs param )
    ]
    [
      copy varName [ letter any [ letter | digit ] ]
    ]
    to end
  ] [ print rejoin [ "ERR: " param ] ]
  varNames: append varNames varName
]
probe varNames
probe optDefs
probe argDefs
probe isArg
probe not isArg
{
    [
      if (not isArg ) [
        "/" copy shortCut letter (print shortCut)
      ]
    ]
}

