package BwCore;
use v5.18;
use strict;
use warnings;
use Exporter;
use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
$VERSION = 1.00;
@ISA = qw(Exporter);
@EXPORT_OK = ();
%EXPORT_TAGS = ();
@EXPORT = qw/
  getFuncName
  camelCaseToKebabCase
  kebabCaseToCamelCase
  hasItem
  shortenFileSpec
  typeOfValue
  validateStruct
  execCmd
  docker
  dockerCompose
  gitIsChangedFile
  mkFileFromTemplate
/;

# =============================================================================

use Data::Dumper;
use BwAnsi;

# =============================================================================

sub getFuncName(;$) {
  return unless defined wantarray;
  my $deep = shift || 0;
  my @caller = caller($deep + 1);
  return unless scalar @caller;
  my @splitted = split '::', $caller[3];
  if (wantarray) {
    @splitted
  } else {
    pop @splitted;
  }
}

sub camelCaseToKebabCase($) {
  $_ = shift;
  s/(?<=.)([A-Z])/-\L$1/g;
  s/^(.)/\L$1/;
  $_;
}

sub kebabCaseToCamelCase($) {
  $_ = shift;
  s/(-)(\w)/\U$2/g;
  $_;
}

sub hasItem($@) {
  my $testItem = shift;
  foreach my $item ( @_ ) {
    return 1 if $item eq $testItem;
  }
  return 0;
}

sub shortenFileSpec($) {
  my $fileSpec = shift || die;
  $fileSpec =~ s|^$ENV{HOME}/|~/|;
  $fileSpec;
}

sub _getAsArrayRef($) {
  my $value = shift;
  my $result = [];
  if ( defined $value ) {
    $result = ref $value eq 'ARRAY' ? $value : [ $value ];
  }
  return $result;
}

sub typeOfValue($) {
  my $value = shift;
  my $result = defined $value ? ref $value : 'undef';
  my %normalizedValueTypes = (
    'HASH' => 'hash',
    'CODE' => 'sub',
    'ARRAY' => 'array',
    'Regexp' => 'regexp',
    '' => 'scalar',
  );
  $normalizedValueTypes{$result} || $result;
}

sub _getSubableValue($$@) {
  my $value = shift;
  my $subParamsDescription = shift or die;
  my $suffix = '';
  if ( typeOfValue $value eq 'sub' ) {
    $value = &{$value};
    $suffix = "->($subParamsDescription)";
  }
  return ($value, $suffix);
}

sub _getValidatedTypeValue {
  my $opt = typeOfValue $_[0] eq 'hash' ? shift : {};
  my $where = shift;
  my $value = shift;
  my @expectedTypes = @_;

  my $valueType = typeOfValue $value;

  my $suffix = '';
  if ($valueType eq 'sub' && exists $opt->{sub} && defined $opt->{sub}) {
    my $subType = typeOfValue $opt->{sub};
    die Dumper({
      err => '_getValidatedTypeValue $opt->{sub} is expected to be an empty string or ref to array with at least two items, where first item is scalar',
      opt => $opt,
    }) unless $subType eq 'scalar' && length $opt->{sub} || $subType eq 'array' && scalar @{$opt->{sub}} >= 2 && typeOfValue $opt->{sub}->[0] eq 'scalar';
    my $subParamsDescription = '';
    if ($subType eq 'scalar') {
      $value = $value->();
    } else {
      my @params = @{$opt->{sub}};
      $subParamsDescription = shift @params;
      $value = $value->(@params);
    }
    $suffix = "->($subParamsDescription)";
  }

  die Dumper({
    where => $where,
    err => 'typeOfValue $value$suffix is not of expectedTypes',
    'typeOfValue $value$suffix' => $valueType,
    expectedTypes => \@expectedTypes,
    "value$suffix" => $value
  }) unless scalar @expectedTypes == 0 || hasItem $valueType, @expectedTypes;

  return !wantarray ? $value : ($value, $valueType);
}

sub _getValidatedTypeValueOfHashKey {
  my $opt = typeOfValue $_[0] eq 'hash' ? shift : {};
  my $whereHash = shift;
  my $hash = _getValidatedTypeValue($whereHash, shift, 'hash');
  my $key = _getValidatedTypeValue('_validateTypeOfHashKeyValue $key', shift, 'scalar');
  my @expectedTypes = @_;

  if (exists $hash->{$key}) {
    _getValidatedTypeValue({ sub => exists $opt->{sub} ? $opt->{sub} : undef }, "$whereHash->{$key}", $hash->{$key}, @expectedTypes);
  } elsif (exists $opt->{isRequired} && $opt->{isRequired}) {
    die Dumper({
      where => $whereHash,
      err => '$key is expected to exist in $hash',
      key => $key,
      hash => $hash,
    });
  } else {
    return !wantarray ? undef : (undef) x 2;
  }
}

sub _checkHashHasOnlyExpectedKeys {
  my $whereHash = shift;
  my $hash = _getValidatedTypeValue($whereHash, shift, 'hash');
  my $expectedKeys = _getValidatedTypeValue('_checkHashHasOnlyExpectedKeys $expectedKeys', shift, 'hash');

  foreach my $key (keys %{$hash}) {
    die Dumper({
      where => $whereHash,
      err => '$hash has non expected $key'
      key => $key,
      expectedKeys => \@expectedKeys,
      hash => $hash,
    }) unless hasItem($key, keys %{$expectedKeys});
  }
}

sub _validateScalar {
  my $where = shift;
  my $value = shift;
  my $struct = shift;
  my $whereStruct = shift;

  my ($enum, $enumType) = _getValidatedTypeValueOfHashKey($whereStruct, $struct, 'enum', 'array');
  if ( $enumType ) {
    die Dumper({
      where => $where,
      err => '$enum has no $value',
      value => $value,
      enum => $struct->{enum}
    })
      unless hasItem($value, @{$struct->{enum}});
  } else {
    my ($nonFalse, $nonFalseType) = _getValidatedTypeValueOfHashKey($whereStruct, $struct, 'noFalse', 'scalar');
    if ($nonFalseType) {
      die Dumper({
        where => $where,
        err => '$value can not be false',
        value => $value,
      }) if $nonFalse && !$value;
    }
  }
}

sub validateStruct {
  my $where = shift or die;
    typeOfValue $where eq 'scalar'
      or die Dumper( { err => '$where is expected to be a scalar', where => $where} );
  my $value = shift;
  my $struct = shift;
  my $whereStruct = shift || '';

  _getValidatedTypeValue($whereStruct, $struct, 'hash');

  my ($type, $typeType) = _getValidatedTypeValueOfHashKey({ isRequired => 1, sub => [ '$struct', $struct ] }, $whereStruct, $struct, 'type', 'scalar', 'array');

  my @types =
    $typeType eq 'scalar' && $type eq 'scalarOrArrayOfScalars' ? ('scalar', 'array') :
    $typeType eq 'array' ? @{$type} : ( $type )
  ;

  my $expectedFields = {
    type => 1,
    validate => 1,
  };
  if (hasItem 'hash', @types) {
    $expectedFields->{keys} = 1;
  }
  if ($type eq 'scalarOrArrayOfScalars' || hasItem 'scalar', $types) {
    $expectedFields->{enum} = 1;
    $expectedFields->{nonFalse} = 1;
  }
  if (hasItem 'array', $types) {
    if (exists $struct->{arrayItem}) {
      $expectedFields->{arrayItem} = 1;
    } else {
      $expectedFields->{value} = 1;
    }
  }
  if (hasItem 'Hash::Ordered', $types) {
    $expectedFields->{value} = 1;
  }
  _checkHashHasOnlyExpectedKeys($whereStruct, $struct, $expectedFields);

  ($value, my $valueType) = _getValidatedTypeValue($where, $value, @types);

  if ( $valueType eq 'hash') {
    my ($keys, $keysType) = _getValidatedTypeValueOfHashKey($whereStruct, $struct, 'keys', 'hash');
    if ($keysType) {
      my %validKeys;
      foreach my $key (keys %{$keys}) {
        my $keyDef = _getValidatedTypeValueOfHashKey({ sub => ['$struct', $struct] }, "$whereStruct\->{keys}", $keys, $key, 'hash');
        next if $keyDef->{condition} && !$keyDef->{condition}->($value);

        my ($keyValue, $keyValueType) = _getValidatedTypeValueOfHashKey({ isRequired => $keyDef->{isRequired} }, "$where\->{$key}", $value, $key);
        if ($keyValueType) {
          $value->{$key} = validateStruct($where . "->{$key}", $value->{$key}, $keyDef);
          $validKeys{$key} = 1;
        }
      }
      _checkHashHasOnlyExpectedKeys($where, $value, \%validKeys);
    }
  } elsif ( $valueType eq 'scalar' ) {
    _validateScalar($where, $value, $struct, $whereStruct);
  } elsif ( $valueType eq 'array' ) {
    if ($type eq 'scalarOrArrayOfScalars') {
      for (my $i = 0; $i < scalar @{$value}; $i++) {
        _validateScalar($where . "->[$i]", $value->[$i], $struct, $whereStruct);
      }
    } else {
      my ($valueStruct, $valueStructType) = _getValidatedTypeValueOfHashKey($whereStruct, $struct, 'arrayItem', 'hash');
      if (!$valueStructType) {
        ($valueStruct, $valueStructType) = _getValidatedTypeValueOfHashKey($whereStruct, $struct, 'value', 'hash');
      }
      die Dumper({
        where => $whereStruct,
        err => 'one of $keys is expected to exist in $struct',
        keys => [ qw/arrayItem value/ ],
        struct => $struct,
      }) unless $valueStructType;

      for (my $i = 0; $i < scalar @{$value}; $i++) {
        $value->[$i] = validateStruct($where . "->[$i]", $value->[$i], $valueStruct);
      }
    }
  } elsif ( $valueType eq 'Hash::Ordered' ) {
    my ($valueStruct, $valueStructType) = _getValidatedTypeValueOfHashKey({ isRequired => 1 }, $whereStruct, $struct, 'value', 'hash');
    foreach my $key ($value->keys) {
      $value->set($key, validateStruct($where . "->get($key)", $value->get($key), $valueStruct));
    }
  } elsif ( ! hasItem($valueType, 'sub', 'undef') ) {
    die Dumper({ err => 'unexpected $valueType', types => $types, '$valueType' => $valueType, value => $value });
  }

  my ($validate, $validateType) = _getValidatedTypeValueOfHashKey($whereStruct, $struct, 'validate', 'sub');
  $value = $validate->($where, $value, $struct, $whereStruct) if $validateType;

  return $value;
}

sub execCmd {
  my $opt = ref $_[0] ne 'HASH' ? {} : validateStruct('execCmd first param', shift, {
    type => 'hash',
    keys => {
      v => {
        type => 'scalar',
        enum => [ 'all', 'allBrief', 'err', 'ok' ],
      },
      return => {
        type => 'scalar',
        enum => [ 'stdout', 'stderr', 'allSeparate', 'allTogether' ],
      },
      silent => {
        type => 'scalar',
        enum => [ 'stdout', 'stderr', 'all' ],
      },
      exitOnError => {
        type => [ 'scalar', 'sub' ],
      },
    },
  });
  my $cmd;
  foreach (@_) {
    my $arg = $_;
    if ( $arg =~ m/[\s"]/ ) {
      $arg =~ s/"/\\"/g;
      $arg = "\"$arg\"";
    };
    $cmd .= " " if length $cmd;
    $cmd .= $arg;
  }
  print ansi "<ansiCmd>$cmd<ansi> . . .\n" if $opt->{v} && $opt->{v} eq 'all';
  my ($stdout, $stderr, $output) = ('') x 2;
  use IPC::Run3;
  run3($cmd, undef,
    sub {
      my $line = shift;
      if ( exists $opt->{return} ) {
        if (hasItem($opt->{return}, qw/stdout allSeparate/)) {
          $stdout .= $line;
        } else {
          $output .= $line;
        }
      }
      print $line unless exists $opt->{silent} && hasItem($opt->{silent}, qw/stdout all/);
    },
    sub {
      my $line = shift;
      if ( exists $opt->{return} ) {
        if (hasItem($opt->{return}, qw/stderr allSeparate/)) {
          $stderr .= $line;
        } else {
          $output .= $line;
        }
      }
      print $line unless exists $opt->{silent} && hasItem($opt->{silent}, qw/stderr all/);
    }
  );
  my $errorCode = ${^CHILD_ERROR_NATIVE} / 256;# https://stackoverflow.com/questions/3736320/executing-shell-script-with-system-returns-256-what-does-that-mean
  my ($ansi, $prefix) = $errorCode == 0 ? ('OK') x 2 : ('Err', 'ERR');
  print ansi $ansi, "$prefix: <ansiCmd>$cmd\n" if $opt->{v} && (
    $opt->{v} =~ /^all/ ||
    $opt->{v} eq 'ok' && $errorCode == 0 ||
    $opt->{v} eq 'err' && $errorCode != 0 ||
  0);
  my $exitOnError = exists $opt->{exitOnError} && $opt->{exitOnError};
  exit $errorCode if $exitOnError && (
    typeOfValue $exitOnError eq 'sub' ? $exitOnError->($errorCode) : $errorCode
  );
  if (!exists $opt->{return}) {
    die Dumper({err => 'wantarray while $opt->{return} not exists', opt => $opt}) if wantarray;
    return $errorCode;
  } else {
    my @result;
    if ($opt->{return} eq 'stdout') {
      push @result, $stdout;
    } elsif ($opt->{return} eq 'stderr') {
      push @result, $stderr;
    } elsif ($opt->{return} eq 'allSeparate') {
      push @result, $stdout, $stderr;
    } else {
      push @result, $output;
    }
    push @result, $errorCode unless $exitOnError;
    if ( scalar @result == 1 ) {
      return wantarray ? @result : $result[0];
    } else {
      die Dumper({err => '!wantarray while $opt->{return}' . ($exitOnError ? '' : ' and no $opt->{exitOnError}') , opt => $opt}) unless wantarray;
    }
    return @result;
  }
}

sub docker {
  my $opt = ref $_[0] eq 'HASH' ? shift : {};
  # TODO: bw_install docker --silentIfAlreadyInstalled || return $?
  unshift @_, 'docker';
  if ( $ENV{OSTYPE} =~ m/^linux/ ) {
    unshift @_, 'sudo';
  } elsif ( $ENV{OSTYPE} !~ m/^darwin/ ) {
    die ansi 'Err', "ERR: Неожиданный тип OS <ansiPrimaryLiteral>$ENV{OSTYPE}";
  }
  unshift @_, $opt;
  &execCmd;
}

sub dockerCompose {
  my $opt = ref $_[0] eq 'HASH' ? shift : {};
  # TODO: bw_install docker-compose --silentIfAlreadyInstalled || return $?
  unshift @_, 'docker-compose';
  if ( $ENV{OSTYPE} =~ m/^linux/ ) {
    unshift @_, 'sudo';
  } elsif ( $ENV{OSTYPE} !~ m/^darwin/ ) {
    die ansi 'Err', "ERR: Неожиданный тип OS <ansiPrimaryLiteral>$ENV{OSTYPE}";
  }
  unshift @_, $opt;
  &execCmd;
}

sub gitIsChangedFile($;$) {
  my $fileName = shift || die;
  my $dir = shift;

  my @command = ! $dir ? () : ('cd', $dir, '&&');
  push @command, qw/git update-index -q --refresh && git diff-index --name-only HEAD --/; #" | grep -Fx \"$fileName\" >/dev/null 2>&1";
  my $stdout = execCmd({ v => 'err', return => 'stdout', silent => 'stdout', exitOnError => 1 }, @command);
  scalar grep { $_ eq $fileName } split /\n/, $stdout;
}

sub mkFileFromTemplate {
  my $opt = ref $_[0] eq 'HASH' ? shift : {};
  my $fileSpec = shift or die;
  my $templateFileSpec = shift or die;
  my $varSubstSub = shift or die;
    ref $varSubstSub eq 'CODE' or die Dumper({varSubstSub => $varSubstSub});

  {
    local $/ = undef;
    open my $fh, $templateFileSpec or die;
    $_ = <$fh>;
    close $fh;
  }

  s(\${([^}]+)})( $varSubstSub->($1) )egm;

  {
    open my $fh, '>', $fileSpec or die;
    unless ( exists $opt->{noNotice} or $opt->{noNotice} ) {
      my ($commentPrefix, $commentSuffix, $commentPreLine, $commentPostLine) = ('') x 4;
      if ($fileSpec =~ m/\.html$/) {
        ($commentPreLine, $commentPostLine) = ( '<!--', '-->' );
      } elsif ($fileSpec =~ m/\.conf$/) {
        $commentPrefix = ('# ');
      } else {
        die Dumper({fileSpec => $fileSpec});
      }
      print $fh <<"HEADER" if $commentPreLine;
${commentPreLine}
HEADER
      print $fh <<"HEADER";
${commentPrefix}ВНИМАНИЕ!!!${commentSuffix}
${commentPrefix}    Этот файл создан автоматически из \"$templateFileSpec\"${commentSuffix}
${commentPrefix}    Поэтому изменения нужно вносить именно туда${commentSuffix}
HEADER
      print $fh <<"HEADER" if $commentPostLine;
${commentPostLine}
HEADER
    }
    print $fh $_;
    close $fh;
  }
}

# =============================================================================
# =============================================================================

1;