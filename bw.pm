package bw;
use strict;
use warnings;

my @vars;
my @export;

BEGIN {
  open(IN, $ENV{_bwDir} . "/bw.pm");
  while (<IN>) {
    if ( /^\s*sub\s+([a-z][\w\d]*)/ ) {
      push @export, "$1";
    } elsif ( /^\$(ansi[\w\d]+)/ ) {
      my $var = "\$$1";
      push @export, $var;
      push @vars, $var;
    }
  }
  close(IN);
}

use Exporter;
use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
use vars @vars;

$VERSION = 1.00;
@ISA = qw(Exporter);

@EXPORT=( @export );
@EXPORT_OK = ();
%EXPORT_TAGS = ();

# =============================================================================

use Hash::Ordered; # https://metacpan.org/pod/Hash::Ordered
use Data::Dumper;
# print Dumper(\@EXPORT);

$ansiReset="\e[0m"; # https://superuser.com/questions/33914/why-doesnt-echo-support-e-escape-when-using-the-e-argument-in-macosx/33950#33950
$ansiBold="\e[1m";
$ansiDim="\e[2m";
$ansiItalic="\e[3m";
$ansiUnderline="\e[4m";
$ansiBlink="\e[5m";
# $ansi6="\e[6m";
$ansiInvert="\e[7m";
$ansiHidden="\e[8m";
$ansiStrike="\e[9m";

$ansiResetBold="\e[22m"; # именно так, а не 20
$ansiResetDim="\e[22m";
$ansiResetItalic="\e[23m";
$ansiResetUnderline="\e[24m";
$ansiResetBlink="\e[25m";
# $ansiReset6="\e[26m";
$ansiResetInvert="\e[27m";
$ansiResetHidden="\e[28m";
$ansiResetStrike="\e[29m";

$ansiBlack="\e[30m";
$ansiRed="\e[31m";
$ansiGreen="\e[32m";
$ansiYellow="\e[33m";
$ansiBlue="\e[34m";
# $ansiMagenta="\e[35m";
$ansiMagenta="\e[38;5;201m"; # https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences/33206814#33206814
$ansiCyan="\e[36m";
$ansiLightGray="\e[37m";
$ansiLightGrey="${ansiLightGray}";
# $ansi38="\e[38m";
$ansiDefault="\e[39m";
$ansiDarkGray="\e[90m";
$ansiDarkGrey="${ansiDarkGray}";
$ansiLightRed="\e[91m";
$ansiLightGreen="\e[92m";
$ansiLightYellow="\e[93m";
$ansiLightBlue="\e[94m";
$ansiLightMagenta="\e[95m";
$ansiLightCyan="\e[96m";
$ansiWhite="\e[97m";

$ansiHeader="${ansiBold}${ansiLightGray}";
$ansiUrl="${ansiBlue}${ansiUnderline}";
$ansiCmd="${ansiWhite}${ansiBold}";
$ansiFileSpec="${ansiWhite}${ansiBold}";
$ansiDir="${ansiWhite}${ansiBold}";
$ansiErr="${ansiRed}${ansiBold}";
$ansiWarn="${ansiYellow}${ansiBold}";
$ansiWill="${ansiYellow}${ansiBold}";
$ansiOK="${ansiGreen}${ansiBold}";
$ansiOutline="${ansiMagenta}${ansiBold}"; # $ansiOutline="${ansiMagenta}${ansiResetBold}"
$ansiDebug="${ansiBlue}${ansiResetBold}";
$ansiPrimaryLiteral="${ansiCyan}${ansiBold}";
$ansiSecondaryLiteral="${ansiCyan}${ansiResetBold}";

sub shortenFileSpec {
  my $fileSpec = shift;
  $fileSpec =~ s|^$ENV{HOME}/|~/|;
  $fileSpec;
}

sub hasItem {
  my $testItem = shift;
  foreach my $item ( @_ ) {
    return 1 if $item eq $testItem;
  }
  return 0;
}

sub getFuncName {
  my $deep = shift || 0;
  my @caller = caller($deep + 1);
  return unless scalar @caller;
  my @splitted = split '::', $caller[3];
  pop @splitted;
}

sub getPackageName {
  my $deep = shift || 0;
  my @caller = caller($deep + 1);
  return unless scalar @caller;
  my @splitted = split '::', $caller[3];
  shift @splitted;
}

sub getSubCommands {
  my $pmFileSpec = shift || die;
  my $deep = shift || 0;
  my @caller = caller($deep + 1);
  return unless scalar @caller;
  my ($packageName, $funcName) = split '::', $caller[3];
  my $glob = $packageName . '::';
  my $result = { 
    byName => Hash::Ordered->new(),
    byNameOrShortcut => Hash::Ordered->new(),
  };
  my $allSubCommands;
  open(IN, $pmFileSpec);
  while (<IN>) {
    next unless /^\s*sub\s+(${funcName}_([\w\d]+))/;
    no strict 'refs';
    my ($funcName, $cmdName) = ($1, camelCaseToKebabCase($2));
    my $def = ${$packageName . "::${funcName}Def"};
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

sub camelCaseToKebabCase {
  $_ = shift;
  s/(?<=.)([A-Z])/-\L$1/g; 
  s/^(.)/\L$1/;
  $_;
}

sub getDescription {
  my $descriptionHolder = shift || die;
  my $descriptionParam = shift;
  my $description = $descriptionHolder->{description};
  $description = $description->($descriptionParam) if ref $description eq 'CODE'; 
  $description;
}

# =============================================================================
# =============================================================================

1;