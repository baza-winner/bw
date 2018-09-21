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
    '' => 'scalar',
  );
  $normalizedValueTypes{$valueType} || $result;
}

sub
sub validateStruct {
  my $where = shift or die;
    typeOfValue $where eq 'scalar'
      or die Dumper( { err => '$where is expected to be a scalar', where => $where} );
  my $value = shift;
  my $struct = shift or die Dumper({ where => $where });
  my $whereStruct = shift || '';

  typeOfValue $struct eq 'hash'
    or die Dumper({ err => 'typeOfValue $struct is expected to be a hash', whereStruct => $whereStruct , 'typeOfValue $struct' => typeOfValue $struct, struct => $struct });

  my @validFields = qw/type/;
  exists $struct->{type}
    or die Dumper({ err => '$struct->{type} is expected to exists', whereStruct => $whereStruct, struct => $struct });
  my $type = $struct->{type};

  my $typeSuffix = '';
  if ( typeOfValue $type eq 'sub' ) {
    $type = $type->($struct);
    $typeSuffix = '->($struct)';
  }
  hasItem(typeOfValue $type, 'scalar', 'array')
    or die Dumper({ err => "\$struct->{type}$typeSuffix is expected to be a sclar or an array",  whereStruct => $whereStruct, "typeOfValue \$struct->{type}$typeSuffix" => typeOfValue $type, type => $type, struct => $struct });

  my $valueType = typeOfValue $value;
  if ( typeOfValue $type eq 'scalar' && $type eq 'scalarOrArrayOfScalars' ) {
    hasItem($valueType, 'scalar', 'array')
      or die Dumper({ err => 'typeOfValue $value is not expected by $struct->{type}', where => $where, 'typeOfValue $value' => $valueType, '$struct->{type}' => $struct->{type}, '$value' => $value });
  } else {
    my $types = typeOfValue($type) eq 'array' ? $type : [ $type ];
    hasItem($valueType, @{$types}) or die Dumper({ err => '$valueType of $value is not expected by $struct->{type}', where => $where, 'typeOfValue $value' => $valueType, '$struct->{type}' => $struct->{type}, '$value' => $value });
  }

  if ( $valueType eq 'hash') {
    if ( $struct->{keys} ) {
      push @validFields, qw/keys/;
      my $keys = $struct->{keys};
      ref $keys eq 'HASH' or die Dumper({ 'ref $struct->{keys}' => ref $keys });
      my @validKeys;
      foreach my $key (keys %{$keys}) {
        my $keyDef = $keys->{$key};

        my $keyDefSuffix = '';
        if ( typeOfValue $keyDef eq 'sub' ) {
          $keyDef = $keyDef->($struct);
          $keyDefSuffix = '->($struct)';
        }
        typeOfValue $keyDef eq 'hash' or die Dumper({ "typeOfValue \$keys->{$key}$keyDefSuffix" => typeOfValue $keyDef });

        next if $keyDef->{condition} && !$keyDef->{condition}->($value);
        !$keyDef->{isRequired} or exists($value->{$key}) or die Dumper({ 'err' => 'required $key is absent in $value', '$key' => $key, '$value' => $value });
        if ( exists($value->{ $key }) ) {
          $value->{$key} = validateStruct("$where\->{$key}", $value->{$key}, $keyDef);
        }
        push @validKeys, $key;
      }
      foreach my $key (keys %{$value}) {
        hasItem($key, @validKeys) or die Dumper({ where => $where, '$key' => $key, '@validKeys' => \@validKeys });
      }
    }
  } elsif ( $valueType eq 'scalar' ) {
    if ( $struct->{enum} ) {
      push @validFields, qw/enum/;
      typeOfValue $struct->{enum} eq 'array' or die Dumper({where => $where, 'ref $struct->{enum}' => ref $struct->{enum}});
      hasItem($value, @{$struct->{enum}}) or die Dumper({ where => $where, err => '$enum has no $value', '$value' => $value, '$enum' => $struct->{enum} });
    }
  } elsif ( $valueType eq 'array' ) {
    if ($type eq 'scalarOrArrayOfScalars') {
      if ( $struct->{enum} ) {
        push @validFields, qw/enum/;
        ref $struct->{enum} eq 'ARRAY' or die Dumper({where => $where, 'ref $struct->{enum}' => ref $struct->{enum}});
        hasItem($value, @{$struct->{enum}}) or die Dumper({ where => $where, err => '$enum has no $value', '$value' => $value, '$enum' => $struct->{enum} });
      }
    } else {
      my $valueStruct;
      if ( exists $struct->{arrayItem} ) {
        push @validFields, qw/arrayItem/;
        $valueStruct = $struct->{arrayItem};
      } elsif ( exists $struct->{value} ) {
        push @validFields, qw/value/;
        $valueStruct = $struct->{value};
      } else {
        die Dumper({ where => $where, err => 'expects arrayItem or value field'});
      }
      if ( $valueStruct ) {
        my $i = 0;
        while ($i < scalar @{$value}) {
          $value->[$i] = validateStruct("$where\->[$i]", $value->[$i], $valueStruct);
          $i += 1;
        }
      }
    }
  } elsif ( $valueType eq 'Hash::Ordered' ) {
    if ( exists $struct->{value} ) {
      push @validFields, qw/value/;
      foreach my $key ($value->keys) {
        $value->set($key, validateStruct("$where\->get($key)", $value->get($key), $struct->{value}));
      }
    }
  } elsif ( ! hasItem($valueType, 'sub', 'undef') ) {
    die Dumper({ err => 'unexpected $valueType', types => $types, '$valueType' => $valueType, value => $value });
  }
  if ( exists $struct->{validate} ) {
    push @validFields, qw/validate/;
    ref $struct->{validate} eq 'CODE' or die Dumper({ where => $where, 'ref $struct->{validate}' => ref $struct->{validate} });
    $value = $struct->{validate}->($where, $value);
  }
  if ( exists $struct->{normalize} ) {
    push @validFields, qw/normalize/;
    $struct->{normalize} eq 'to array' or die Dumper({ where => $where, '$struct->{normalize}' => $struct->{normalize} });
    if ( $valueType ne 'array' ) {
      $value = [ $value ];
    }
  }
  return $value;
}

sub execCmd {
  # my $validateStdStream = sub {
  #   my ($where, $value) = @_;
  #   if (typeOfValue($value) eq 'scalar') {


  #   }
  #   return $value;
  # };
  my $opt = ref $_[0] ne 'HASH' ? {} : _validateStruct('execCmd first param', shift, {
    type => 'hash',
    keys => {
      v => {
        type => 'scalar',
        enum => [ 'all', 'allBrief', 'err', 'ok' ],
      },
      stdout => {
        type => [ 'scalar', 'array'],
        enum => [],
        arrayItem => {
          type => 'scalar'
        },
        validate =>
        normalize => 'to array',
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
  my ($stdout, $stderr) = ('') x 2;
  use IPC::Run3;
  run3($cmd, undef,
    sub {
      my $line = shift;
      $stdout .= $line;
      print $line;
    },
    sub {
      my $line = shift;
      $stderr .= $line;
      print $line;
    }
  );
  # if (hasItem 'stdout', @{_getAsArrayRef($opt->{return})}) {
  #   $stdout = qx/$cmd/;
  # } else {
  #   system($cmd);
  # }
  my $errorCode = ${^CHILD_ERROR_NATIVE} / 256;# https://stackoverflow.com/questions/3736320/executing-shell-script-with-system-returns-256-what-does-that-mean
  # print Dumper($errorCode);
  my ($ansi, $prefix) = $errorCode == 0 ? ('OK') x 2 : ('Err', 'ERR');
  print ansi $ansi, "$prefix: <ansiCmd>$cmd\n" if $opt->{v} && (
    $opt->{v} =~ /^all/ ||
    $opt->{v} eq 'ok' && $errorCode == 0 ||
    $opt->{v} eq 'err' && $errorCode != 0 ||
  0);
  my @optReturn = @{_getAsArrayRef($opt->{return})};
  if ( !wantarray ) {
    my $item = shift @optReturn;
    if (defined $item) {
      if ($item eq 'stdout') {
        return $stdout;
      } elsif ($item eq 'stderr') {
        return $stderr;
      }
    }
    return $errorCode;
  } else {
    my @result = ();
    foreach my $item (@optReturn) {
      if (defined $item) {
        if ($item eq 'stdout') {
          push @result, $stdout;
        } elsif ($item eq 'stderr') {
          push @result, $stderr;
        }
      }
    }
    push @result, $errorCode;
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
  my ($stdout, $errorCode) = execCmd({ return => 'stdout', v => 'err' }, @command); exit $errorCode if $errorCode;
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