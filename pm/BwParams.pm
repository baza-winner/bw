package BwParams;
use v5.18;
use strict;
use warnings;
use Exporter;
use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
$VERSION = 1.00;
@ISA = qw(Exporter);
@EXPORT_OK = ();
%EXPORT_TAGS = ();
@EXPORT = ( qw/run processParams preprocessDefs/);

# =============================================================================

use BwAnsi;
use BwCore;
use Data::Dumper;

# =============================================================================

sub run {
  my $packageName = shift or die;
  my $entry = shift or die;
  my $cnf = shift or die;
  my @argv = @_;

  my $fileSpec = $packageName;
  $fileSpec =~ s/::/\//g;
  $fileSpec = "$ENV{_bwDir}/pm/$fileSpec.pm";
  require "$fileSpec";
  my $entityName = 'defs';
  my @splitted = split '::', $packageName;
  $packageName = pop @splitted;
  no strict 'refs';
  my $defs = ${"$packageName::${entityName}"} or die;

  my $def = $defs->{$entry} or die;

  $entityName = 'preprocessCnf';
  if ( exists ${"${packageName}::"}{$entityName} ) {
    my $preprocessCnf = \&{"$packageName::${entityName}"};
    $cnf = $preprocessCnf->($entry, $cnf);
  }
  &{"$packageName::${entry}"}($def, $cnf, @argv);
}

sub _getDescription($$$) {
  my $descriptionContainer = shift or die;
  my $cnf = shift or die;
  my $deep = shift or die;
  my $description = $descriptionContainer->{description} or die;
  ansi ( ref $description ne 'CODE' ? $description : $description->($cnf));
}

sub _printHelp {
  my $def = shift or die;
  my $cnf = shift or die;
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
      my $optDef = $def->{options}->get($optName);
      my $optTitleSuffix = $optDef->{type} eq 'bool' ? '' : ' <ansiOutline>значение<ansi>';
      my $optTitle = "<ansiCmd>--$optName<ansi>$optTitleSuffix";
      if ( $optDef->{shortcut} ) {
        foreach my $shortcut (@{$optDef->{shortcut}}) {
          $optTitle .= " или <ansiCmd>-$shortcut<ansi>$optTitleSuffix";
        }
      }
      print ansi <<"HELP";
    $optTitle
        ${\_getDescription($optDef, $cnf, 2)}
HELP
      my $typeDescription;
      if ($optDef->{type} eq 'int') {
        $typeDescription = 'целое число';
        if ( exists $optDef->{min} ) {
          if ( exists $optDef->{max} ) {
            $typeDescription .= " из диапазона <ansiSecondaryLiteral>$optDef->{min}<ansi>..<ansiSecondaryLiteral>$optDef->{max}<ansi>";
          } else {
            $typeDescription .= " не менее <ansiSecondaryLiteral>$optDef->{min}<ansi>";
          }
        } elsif ( exists $optDef->{max} ) {
          $typeDescription .= " не более <ansiSecondaryLiteral>$optDef->{max}<ansi>";
        }
      }
      if ($typeDescription) {
        print ansi <<"HELP";
        <ansiOutline>Значение<ansi> - $typeDescription
HELP
      }
      if ( exists $optDef->{default} ) {
        print ansi <<"HELP";
        <ansiOutline>Значение<ansi> по умолчанию: <ansiPrimaryLiteral>$optDef->{default}
HELP
      }
    }
  }
  print ansi <<"HELP";
    <ansiCmd>--help<ansi> или <ansiCmd>-?<ansi> или <ansiCmd>-h<ansi>
        Выводит справку
HELP
}

sub _validateOptionDef {
  my $optDef = _validateStruct('_validateOptionDef arg', shift, { type => 'hash' });
  # TODO
  return $optDef;
}

sub _getEntity {
  my $def = shift or die;
  my $cnf = shift or die;
  my $entityName = shift or die;
  my $entity = _validateStruct("\$def->{$entityName}", $def->{$entityName}, {
    type => [ 'Hash::Ordered', 'undef' ],
  }) or return;

  if ($entityName eq 'options' && $cnf->{mixin}) {
    my $funcName = getFuncName(2);
    if ($cnf->{mixin}->{$funcName} && $cnf->{mixin}->{$funcName}->{options}) {
      my $mixinOptions = _validateStruct(
        "\$cnf->{mixin}->{$funcName}->{options}",
        $cnf->{mixin}->{$funcName}->{options}, 
        {
          type => 'Hash::Ordered',
          value => {
            type => 'hash',
            keys => {
              default => {
                type => 'scalar',
              },
            },
          },
        },
      );
      foreach my $key ($mixinOptions->keys) {
        $entity->exists($key) or die Dumper({ err => '$key does not exist in $entity', '$key' => $key, '$entity' => $entity });
        my $value = $entity->get($key);
        my $mixin = $mixinOptions->get($key);
        @{$value}{keys %{$mixin}} = values %{$mixin};
        $entity->set($key, _validateOptionDef($value, $key, $funcName));
      }
    }
  }
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
  $result;
}

sub processParams {
  my $def = shift or die;
  my $cnf = shift or die;

  my $subCommands = _getEntity($def, $cnf, 'subCommands');
  my $options = _getEntity($def, $cnf, 'options');
  # my $args = _getEntity($def, $cnf, 'args');
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
  my $where = shift or die;
    ref $where eq '' or die Dumper( { err => '$where is expected to be a scalar', '$where' => $where} );
  my $value = shift;
  my $struct = shift or die Dumper({ where => $where });
  ref $struct eq 'HASH' or die Dumper({ where => $where , 'ref $struct' => ref $struct });
  my $type = $struct->{type} or die Dumper({ '$struct' => $struct });
  hasItem(ref $type, '', 'ARRAY') or die Dumper({ 'ref $type' => ref $type });
  my $types = ref $type eq 'ARRAY' ? $type : [ $type ];
  my $valueType = defined $value ? ref $value : 'undef';
  my %normalizedValueTypes = (
    'HASH' => 'hash',
    'CODE' => 'sub',
    'ARRAY' => 'array',
    '' => 'scalar',
  );
  my $normalizedValueType = $normalizedValueTypes{$valueType} || $valueType;
  hasItem($normalizedValueType, @{$types}) or die Dumper({ err => '$normalizedValueType of $value is not expected by $struct->{type}', where => $where, '$normalizedValueType' => $normalizedValueType, '$struct->{type}' => $struct->{type}, '$value' => $value });
  if ( $normalizedValueType eq 'hash') {
    if ( $struct->{keys} ) {
      my $keys = $struct->{keys};
      ref $keys eq 'HASH' or die Dumper({ 'ref $struct->{keys}' => ref $keys });
      my @validKeys;
      foreach my $key (keys %{$keys}) {
        my $keyDef = $keys->{$key};
        ref $keyDef eq 'HASH' or die Dumper({ "ref \$keys->{$key}" => ref $keyDef });
        next if $keyDef->{condition} && !$keyDef->{condition}->($value);
        !$keyDef->{isRequired} or exists($value->{$key}) or die Dumper({ 'err' => 'required $key is absent in $value', '$key' => $key, '$value' => $value });
        if ( exists($value->{ $key }) ) {
          $value->{$key} = _validateStruct("$where\->{$key}", $value->{$key}, $keyDef);
        }
        push @validKeys, $key;
      }
      foreach my $key (keys %{$value}) {
        hasItem($key, @validKeys) or die Dumper({ where => $where, '$key' => $key, '@validKeys' => \@validKeys });
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
      ref $struct->{enum} eq 'ARRAY' or die Dumper({where => $where, 'ref $struct->{enum}' => ref $struct->{enum}});
      hasItem($value, @{$struct->{enum}}) or die Dumper({ where => $where, err => '$enum has no $value', '$value' => $value, '$enum' => $struct->{enum} });
    }
  } elsif ( ! hasItem($normalizedValueType, 'sub', 'undef') ) {
    die Dumper({ err => 'unexpected $valueType', types => $types, '$valueType' => $valueType, value => $value });
  }
  if ( $struct->{validate} ) {
    ref $struct->{validate} eq 'CODE' or die Dumper({ where => $where, 'ref $struct->{validate}' => ref $struct->{validate} });
    $value = $struct->{validate}->($where, $value);
  }
  if ( $struct->{normalize} ) {
    $struct->{normalize} eq 'to array' or die Dumper({ where => $where, '$struct->{normalize}' => $struct->{normalize} });
    if ( $normalizedValueType ne 'array' ) {
      $value = [ $value ];
    }
  }
  return $value;
}

sub _preprocessDef {
  my $packageName = shift or die;
  my $allDefs = shift or die;
  my $funcName = shift or die;
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
      length $value == 1 or die Dumper({ where => $where, '$value' => $value,  'length $value' => length $value});
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
                enum => [ 'bool', 'int', 'scalar', 'list' ],
              },
              min => {
                condition => sub { $_[0]->{type} eq 'int' },
                type => 'scalar',
              },
              max => {
                condition => sub { $_[0]->{type} eq 'int' },
                type => 'scalar',
              },
              itemType => {
                condition => sub { $_[0]->{type} eq 'list' },
                type => 'scalar',
                enum => [ 'enum', 'int' ],
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

  if ($def->{options}) {
    ref $def->{options} eq 'Hash::Ordered' or die Dumper({ 'ref $def->{options}' => ref $def->{options} });
    my $options = Hash::Ordered->new();
    foreach my $key ($def->{options}->keys) {
      my $value = $def->{options}->get($key);
      my $key = camelCaseToKebabCase($key);
      $options->set($key, $value);
    }
    $def->{options} = $options;
  }

  $def;
}

sub preprocessDefs {
  my @caller = caller(1);
  $caller[6] =~ m/([\w\d]+)\.pm$/ or die;
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