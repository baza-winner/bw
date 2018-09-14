package BwCore;
use v5.18;
use strict;
use warnings;

my $selfFileSpec;
BEGIN {
  use Exporter;
  use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
  $VERSION = 1.00;
  @ISA = qw(Exporter);
  @EXPORT = ( );
  @EXPORT_OK = ();
  %EXPORT_TAGS = ();

  my @caller = caller(0);
  $selfFileSpec = $caller[1];
  open(IN, $selfFileSpec);
  while (<IN>) {
    push @EXPORT, $1 if /^\s*sub\s+([a-z][\w\d]*)/;
  }
  close(IN);
}

# =============================================================================

sub run {
  my $packageName = shift || die;
  my $entry = shift || die;
  my $cnf = shift || die;
  my @argv = @_;

  my $fileSpec = $packageName;
  $fileSpec =~ s/::/\//g;
  $fileSpec = "$ENV{_bwDir}/pm/$fileSpec.pm";
  require "$fileSpec";
  my $entityName = 'defs';
  my @splitted = split '::', $packageName;
  $packageName = pop @splitted;
  no strict 'refs';
  my $defs = ${"$packageName::${entityName}"} || die;

  my $def = $defs->{$entry} || die;

  $entityName = 'preprocessCnf';
  if ( exists ${"${packageName}::"}{$entityName} ) {
    my $preprocessCnf = \&{"$packageName::${entityName}"};
    $cnf = $preprocessCnf->($entry, $cnf);
  }
  &{"$packageName::${entry}"}($def, $cnf, @argv);
}

sub _getDescription($$$) {
  my $descriptionContainer = shift || die;
  my $cnf = shift || die;
  my $deep = shift || die;
  my $description = $descriptionContainer->{description} || die;
  ansi ( ref $description ne 'CODE' ? $description : $description->($cnf));
}

sub _printHelp {
  my $def = shift || die;
  my $cnf = shift || die;
  my $subCommands = shift;

  my $optionsTitle = $def->{options} ? 'Опции' : 'Опция';
  my $argsTitle = $subCommands ? ' <ansiOutline>Команда<ansi>' : '';
  print ansi <<"HELP";
<ansiHeader>Использование:<ansi> <ansiCmd>$cnf->{cmd}<ansi> [<ansiOutline>$optionsTitle<ansi>]$argsTitle
<ansiHeader>Описание:<ansi> ${\_getDescription($def, $cnf, 1)}
HELP
  if ( $subCommands ) {
    print ansi <<"HELP";
Команда - один из следующих вариантов: <ansiSecondaryLiteral>@{[ join(" ", $subCommands->{byNameOrShortcut}->keys)  ]}<ansiReset>
HELP
    foreach my $cmdName ($subCommands->{byName}->keys) {
      my $cmd = $subCommands->{byName}->get($cmdName);
      my @nameAndShortcuts;
      push @nameAndShortcuts, @{$cmd->{def}->{shortcut}} if $cmd->{def}->{shortcut};
      push @nameAndShortcuts, $cmdName;
      foreach my $cmdNameOrShortcut (@nameAndShortcuts) {
        print ansi <<"HELP";
    <ansiCmd>$cnf->{cmd} $cmdNameOrShortcut<ansi>
HELP
      }
      print ansi <<"HELP";
      ${\_getDescription($cmd->{def}, $cnf, 2)}
HELP
    }
    print ansi <<"HELP";
Подробнее см.
    <ansiCmd>$cnf->{cmd} <ansiOutline>Команда <ansiCmd>--help<ansi>
    <ansiCmd>$cnf->{cmd} <ansiOutline>Команда <ansiCmd>-?<ansi>
    <ansiCmd>$cnf->{cmd} <ansiOutline>Команда <ansiCmd>-h<ansi>
HELP
  }
  print ansi <<"HELP";
$optionsTitle
HELP
  if ( $def->{options} ) {
    for my $optName ($def->{options}->keys) {
      my $optTitle = "<ansiCmd>--$optName<ansi>";
      my $optDef = $def->{options}->get($optName);
      if ( $optDef->{shortcut} ) {
        foreach my $shortcut (@{$optDef->{shortcut}}) {
          $optTitle .= " или <ansiCmd>-$shortcut<ansi>";
        }
      }
      print ansi <<"HELP";
    $optTitle
        ${\_getDescription($optDef, $cnf, 2)}
HELP
    }
  }
  print ansi <<"HELP";
    <ansiCmd>--help<ansi> или <ansiCmd>-?<ansi> или <ansiCmd>-h<ansi>
        Выводит справку
HELP
}

sub _getEntity {
  my $def = shift || die;
  my $cnf = shift || die;
  my $entityName = shift || die;
  my $entity = $def->{$entityName} || return;
  ref $entity eq 'Hash::Ordered' || die Dumper({ 'ref $entity' => ref $entity });

  my $result = {
    byName => Hash::Ordered->new(),
    byNameOrShortcut => Hash::Ordered->new(),
    byShortcut => Hash::Ordered->new(),
  };
  my $all;
  for my $name ($entity->keys) {
    no strict 'refs';
    my $value = $entity->get($name);
    my $def = $value->{def};
    next if $def->{condition} && !$def->{condition}->($cnf);
    $value->{name} = $name;
    $result->{byName}->set($name => $value);
    $result->{byNameOrShortcut}->set($name => $value);
    my $shortcuts = $value->{shortcut};
    if ( $shortcuts ) {
      foreach my $shortcut ( @{$shortcuts} ) {
        $result->{byNameOrShortcut}->set($shortcut => $value);
        $result->{byShortcut}->set($shortcut => $value);
      }
    }
  }
  # $result->{nameAndShortcuts} = join " ", $result->{byNameOrShortcut}->keys;
  # $result->{shortcuts} = join " ", $result->{byShortcut}->keys;
  $result;
}

sub processParams {
  my $def = shift || die;
  my $cnf = shift || die;

  my $subCommands = _getEntity($def, $cnf, 'subCommands');
  my $options = _getEntity($def, $cnf, 'options');
  my $args = _getEntity($def, $cnf, 'args');
  my $result = {};

  my $param = shift;
  while (defined $param) {
    my $optionDef;
    if ($param eq '-?' || $param eq '-h' || $param eq '--help') {
      _printHelp($def, $cnf, $subCommands);
      return undef;
    } elsif ( $param =~ m/^--(.*)/ ) {
      my $optionName = $1;
      if ( !$options || !($optionDef = $options->{byName}->get($optionName))) {
        die ansi 'Err', "<ansiCmd>$cnf->{cmd}<ansi> не ожидает опцию <ansiCmd>--$optionName";
      } elsif ( $optionDef->{type} eq 'bool' ) {
        $result->{$optionName}->{value} = 1;
        $result->{$optionName}->{asis} = [ $param ];
      } else {
        die Dumper({ _ => 'TODO', optionDef => $optionDef });
      }
    } elsif ( $param =~ m/^-(.*)/ ) {
      foreach (split //, $1) {
        if ( !$options || !($optionDef = $options->{byShortcut}->get($_))) {
          die ansi 'Err', "<ansiCmd>$cnf->{cmd}<ansi> не ожидает опцию <ansiCmd>-$_";
        } else {
          my $optionName = $optionDef->{name};
          if ( $optionDef->{type} eq 'bool' ) {
            $result->{$optionName}->{value} = 1;
            $result->{$optionName}->{asis} = [ $param ];
          }
        }
      }
    } elsif ( $subCommands ) {
      my $subCommand = !$param ? undef : $subCommands->{byNameOrShortcut}->get($param);
      if ( $subCommand ) {
        $cnf->{cmd} .= " $param";
        no strict 'refs';
        &{$subCommand->{funcName}}($subCommand->{def}, $cnf, @_);
        return undef;
      } else {
        my $prefix = !$param ? 'в качесте первого аргумента' : "вместо <ansiPrimaryLiteral>$param<ansi>";
        die ansi 'Err', <<"MSG";
ERR: <ansiCmd>$cnf->{cmd}<ansi> ${prefix} ожидает одну из следующих команд <ansiSecondaryLiteral>@{[ join(" ", $subCommands->{byNameOrShortcut}->keys) ]}
MSG
      }
    } else {

    }
    $param = shift;
  }
  return wantarray ? ($result, $cnf) : $result;
}

sub _validateStruct {
  my $where = shift || die;
  my $value = shift;
  my $struct = shift || die Dumper({ where => $where });
  ref $struct eq 'HASH' || die Dumper({ where => $where , 'ref $struct' => ref $struct });
  my $type = $struct->{type} || die Dumper({ '$struct' => $struct });
  hasItem(ref $type, '', 'ARRAY') || die Dumper({ 'ref $type' => ref $type });
  my $types = ref $type eq 'ARRAY' ? $type : [ $type ];
  my $valueType = ref $value;
  my %normalizedValueTypes = (
    'HASH' => 'hash',
    'CODE' => 'sub',
    'ARRAY' => 'array',
    '' => 'scalar',
  );
  my $normalizedValueType = $normalizedValueTypes{$valueType} || $valueType;
  hasItem($normalizedValueType, @{$types}) || die Dumper({ where => $where, '$normalizedValueType' => $normalizedValueType, '$struct->{type}' => $struct->{type} });
  if ( $normalizedValueType eq 'hash') {
    if ( $struct->{keys} ) {
      my $keys = $struct->{keys};
      ref $keys eq 'HASH' || die Dumper({ 'ref $struct->{keys}' => ref $keys });
      my @validKeys;
      foreach my $key (keys %{$keys}) {
        my $keyDef = $keys->{$key};
        ref $keyDef eq 'HASH' || die Dumper({ "ref \$keys->{$key}" => ref $keyDef });
        !$keyDef->{isRequired} || exists $value->{$key} || die Dumper({ '$key' => $key, '$value' => $value });
        if (exists $value->{$key}) {
          $value->{$key} = _validateStruct("$where\->{$key}", $value->{$key}, $keyDef);
        }
        push @validKeys, $key;
      }
      foreach my $key (keys %{$value}) {
        hasItem($key, @validKeys) || die Dumper({ where => $where, '$key' => $key, '@validKeys' => \@validKeys });
      }
    }
  } elsif ( $normalizedValueType eq 'array' ) {
    my $valueStruct = $struct->{arrayItem} || $struct->{value};
    if ( $valueStruct ) {
      my $i = 0;
      while ($i < scalar @{$value}) {
        $value->[$i] = _validateStruct("$where\->[$i]", $value->[$i], $valueStruct);
      }
    }
  } elsif ( $normalizedValueType eq 'Hash::Ordered' ) {
    if ( $struct->{value} ) {
      foreach my $key ($value->keys) {
        $value->set($key, _validateStruct("$where\->get($key)", $value->get($key), $struct->{value}));
      }
    }
  } elsif ( $normalizedValueType eq 'scalar' ) {
    if ( $struct->{enum} ) {
      ref $struct->{enum} eq 'ARRAY' || die Dumper({where => $where, 'ref $struct->{enum}' => ref $struct->{enum}});
      hasItem($value, @{$struct->{enum}}) || die Dumper({ where => $where, value => $value, enum => $struct->{enum} });
    }
  } elsif ( ! hasItem($normalizedValueType, 'sub') ) {

    die Dumper({ _ => 'TODO', types => $types, valueType => $valueType, value => $value });
  }
  if ( $struct->{validate} ) {
    ref $struct->{validate} eq 'CODE' || die Dumper({ where => $where, 'ref $struct->{validate}' => ref $struct->{validate} });
    $value = $struct->{validate}->($where, $value);
  }
  if ( $struct->{normalize} ) {
    $struct->{normalize} eq 'to array' || die Dumper({ where => $where, '$struct->{normalize}' => $struct->{normalize} });
    if ( $normalizedValueType ne 'array' ) {
      $value = [ $value ];
    }
  }
  return $value;
}

sub _preprocessDef {
  my $packageName = shift || die;
  my $allDefs = shift || die;
  my $funcName = shift || die;
  my $validateCmdShortcut = sub {
    my ($where, $value) = @_;
    if (ref $value eq '') {
      $value = camelCaseToKebabCase($value);
    }
    $value;
  };
  my $validateOptionShortcut = sub {
    my ($where, $value) = @_;
    if (ref $value eq '') {
      length $value == 1 || die Dumper({ where => $where, '$value' => $value,  'length $value' => length $value});
    }
    $value;
  };
  my $def = _validateStruct("\$allDefs->{$funcName}", $allDefs->{$funcName},
    {
      type => 'hash',
      keys => {
        condition => {
          type => 'sub',
        },
        description => {
          isRequired => 1,
          type => [ 'scalar', 'sub' ],
        },
        isWrapper => {
          type => 'scalar',
        },
        shortcut => {
          type => [ 'scalar', 'array' ],
          value => {
            type => 'scalar',
            validate => $validateCmdShortcut,
          },
          validate => $validateCmdShortcut,
          normalize => 'to array',
        },
        options => {
          type => 'Hash::Ordered',
          value => {
            type => 'hash',
            keys => {
              type => {
                isRequired => 1,
                type => 'scalar',
                enum => [ 'bool', 'scalar', 'list' ],
              },
              shortcut => {
                type => [ 'scalar', 'array' ],
                value => {
                  type => 'scalar',
                  validate => $validateOptionShortcut,
                },
                validate => $validateOptionShortcut,
                normalize => 'to array',
              },
              description => {
                isRequired => 1,
                type => [ 'scalar', 'sub' ],
              },
            },
          },
        },
        args => {
          type => 'Hash::Ordered',
          value => {
            type => 'hash',
            keys => {
              type => {
                isRequired => 1,
                type => 'scalar',
              },
              description => {
                isRequired => 1,
                type => [ 'scalar', 'sub' ],
              },
            },
          },
        },
      },
    }
  );

  if ($def->{isWrapper}) {

    my $glob = $packageName . '::';
    my $subCommands = Hash::Ordered->new();
    foreach (keys %{$allDefs}) {
      next unless /^(${funcName}_([a-zA-Z\d]+))(?![\w\d])/;
      no strict 'refs';
      my ($funcName, $cmdName) = ($1, camelCaseToKebabCase($2));
      my $def = ${ "${packageName}::${funcName}Def" };
      my $value = {
        funcName => "${packageName}::${funcName}",
        def => _preprocessDef($packageName, $allDefs, $funcName),
      };
      $subCommands->set($cmdName => $value);
    }

    $def->{subCommands} = $subCommands;
  }

  $def;
}

sub preprocessDefs {
  my @caller = caller(1);
  $caller[6] =~ m/([\w\d]+)\.pm$/ || die;
  my $packageName = $1;
  my $varName = "selfFileSpec";
  no strict 'refs';
  my $pmFileSpec = ${"$packageName::${varName}"}; # без подстановки varName не работает
  my $allDefs = {};

  open(IN, $pmFileSpec);
  while (<IN>) {
    next unless /^\s*sub\s+((\w[\w\d]*))/;
    my $funcName = $1;
    my $defVarName = "${funcName}Def";
    my $defPackageVarName = "$packageName::$defVarName";
    next unless ${$packageName . "::"}{$defVarName};
    my $def = ${$defPackageVarName};
    next unless $def;
    $allDefs->{$funcName} = $def;
  }

  my $defs = {};
  foreach (keys %{$allDefs}) {
    next unless /^(([[a-zA-Z][a-zA-Z\d]*))(?![\w\d])/;
    my $funcName = $1;
    $defs->{$funcName} = _preprocessDef($packageName, $allDefs, $funcName);
  }

  return $defs;
}

# =============================================================================
# =============================================================================

1;