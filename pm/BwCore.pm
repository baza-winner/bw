package BwCore;
use v5.18;
use strict;
use warnings;

my $selfFileSpec;
BEGIN {
  use Hash::Ordered; # https://metacpan.org/pod/Hash::Ordered
  use Data::Dumper;

  use File::Find qw/find/;
  my $carpAlwaysIsInstalled;
  no warnings 'File::Find';
  find { wanted => sub { $carpAlwaysIsInstalled = 1 if /\/Carp\/Always(?:\.pm)?$/ }, no_chdir => 1 }, @INC;
  if ( $carpAlwaysIsInstalled ) {
    require Carp::Always;# https://metacpan.org/pod/Carp::Always
    Carp::Always->import;
  }

  use BwAnsi;

  use Exporter;
  use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
  $VERSION = 1.00;
  @ISA = qw(Exporter);
  @EXPORT = ( 'ansi', 'Dumper',  );
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

sub shortenFileSpec($) {
  my $fileSpec = shift || die;
  $fileSpec =~ s|^$ENV{HOME}/|~/|;
  $fileSpec;
}

sub hasItem($@) {
  my $testItem = shift;
  foreach my $item ( @_ ) {
    return 1 if $item eq $testItem;
  }
  return 0;
}

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

# =============================================================================

sub _getSubCommands($$$$) {
  my ($pmFileSpec, $packageName, $funcName, $cnf)  = @_;

  my $glob = $packageName . '::';
  my $result = {
    byName => Hash::Ordered->new(),
    byNameOrShortcut => Hash::Ordered->new(),
  };
  my $allSubCommands;
  open(IN, $pmFileSpec);
  while (<IN>) {
    next unless /^\s*sub\s+(${funcName}_([a-zA-Z\d]+))/;
    no strict 'refs';
    my ($funcName, $cmdName) = ($1, camelCaseToKebabCase($2));
    my $def = ${$packageName . "::${funcName}Def"};
    # print Dumper( { def => $def, cnf => $cnf } );
    next if $def->{condition} && !$def->{condition}->($cnf);
    my $value = {
      funcName => $funcName,
      def => $def,
    };
    $result->{byName}->set($cmdName => $value);
    $result->{byNameOrShortcut}->set($cmdName => $value);
    if ( $def->{shortcuts} ) {
      @{$def->{shortcuts}} = map { camelCaseToKebabCase($_) } @{$def->{shortcuts}};
      foreach my $shortcut ( @{$def->{shortcuts}} ) {
        $result->{byNameOrShortcut}->set($shortcut => $value);
      }
    }
  }
  close(IN);
  $result->{cmdNameAndShortcuts} = join " ", $result->{byNameOrShortcut}->keys;
  $result;
}

sub _getDescription($$) {
  my ($descriptionContainer, $cnf) = @_;
  my $description = $descriptionContainer->{description} || die;
  # print Dumper({description => $description, descriptionContainer => $descriptionContainer});
  ansi ( ref $description ne 'CODE' ? $description : $description->($cnf));
}

sub _printHelp($$$) {
  my ($cnf, $def, $subCommands) = @_;
  # print Dumper({def => $def});
  my $optionsTitle = $def->{options} ? 'Опции' : 'Опция';
  my $argsTitle = $subCommands ? ' <ansiOutline>Команда<ansi>' : '';
  print ansi <<"HELP";
<ansiHeader>Использование:<ansi> <ansiCmd>$cnf->{cmd}<ansi> [<ansiOutline>$optionsTitle<ansi>]$argsTitle
<ansiHeader>Описание:<ansi> ${\_getDescription($def, $cnf)}
HELP
  if ( $subCommands ) {
    print ansi <<"HELP";
Команда - один из следующих вариантов: <ansiSecondaryLiteral>$subCommands->{cmdNameAndShortcuts}<ansiReset>
HELP
    foreach my $cmdName ($subCommands->{byName}->keys) {
      my $cmd = $subCommands->{byName}->get($cmdName);
      my @cmdNameAndShortcuts;
      push @cmdNameAndShortcuts, @{$cmd->{def}->{shortcuts}} if $cmd->{def}->{shortcuts};
      push @cmdNameAndShortcuts, $cmdName;
      foreach my $cmdNameOrShortcut (@cmdNameAndShortcuts) {
        print ansi <<"HELP";
    <ansiCmd>$cnf->{cmd} $cmdNameOrShortcut<ansi>
HELP
      }
      print ansi <<"HELP";
      ${\_getDescription($cmd->{def}, $cnf)}
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
  print ansi <<"HELP";
    <ansiCmd>--help<ansi> или <ansiCmd>-?<ansi> или <ansiCmd>-h<ansi>
        Выводит справку
HELP
}

sub _processWrapper($$$$@) {
  my ($cnf, $packageName, $subCommands, $param, @params) = @_;
  # print Dumper($subCommands);
  my $subCommand = ! $param ? undef : $subCommands->{byNameOrShortcut}->get($param);
  if ( $subCommand ) {
    $cnf->{cmd} .= " $param";
    # print Dumper({ cnf => $cnf });
    &{\&{"$packageName::$subCommand->{funcName}"}}($cnf, @params); # https://stackoverflow.com/questions/1915616/how-can-i-elegantly-call-a-perl-subroutine-whose-name-is-held-in-a-variable
  } else {
    my $prefix = !$param ? 'в качесте первого аргумента' : "вместо <ansiPrimaryLiteral>$param<ansi>";
    die ansi 'Err', "ERR: <ansiCmd>$cnf->{cmd}<ansi> ${prefix} ожидает одну из следующих команд <ansiSecondaryLiteral>$subCommands->{cmdNameAndShortcuts}";
  }
}

sub processParams($@) {
  my $cnf = shift;
  # print Dumper($param);

  no strict 'refs';
  my ($packageName, $funcName) = getFuncName 1;
  my $def = ${"$packageName::${funcName}Def"};

  my $subCommands;
  if ( $def->{isWrapper} ) {
    my $varName = "selfFileSpec";
    my $packageFileSpec = ${"$packageName::${varName}"}; # без подстановки varName не работает
    $subCommands = _getSubCommands $packageFileSpec, $packageName, $funcName, $cnf;
  }

  my $param = shift;
  # if $def->{options}

  if ($param and ( $param eq '-?' || $param eq '-h' || $param eq '--help')) {
    _printHelp $cnf, $def, $subCommands;
  } else {
    if ( $def->{isWrapper} ) {
      _processWrapper $cnf, $packageName, $subCommands, $param, @_;
    }
  }
}

# =============================================================================
# =============================================================================

1;