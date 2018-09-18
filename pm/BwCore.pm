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
  docker 
  execCmd 
  gitIsChangedFile 
  shortenFileSpec
  hasItem
  getFuncName
  camelCaseToKebabCase
/;

# =============================================================================

use BwAnsi;

# =============================================================================

sub docker {
  my $opt = {};
  $opt = shift if ref $_[0] eq 'HASH';
  # TODO: bw_install docker --silentIfAlreadyInstalled || return $?
  unshift @_, 'docker';
  if ( $ENV{OSTYPE} =~ m/^linux/ ) {
    unshift @_, 'sudo';
  } elsif ( $ENV{OSTYPE} !~ m/^darwin/ ) {
    die ansi 'Err', "ERR: Неожиданный тип OS <ansiPrimaryLiteral>$ENV{OSTYPE}";
  }
  unshift @_, $opt;
  !&execCmd;
}

sub execCmd {
  my $opt = {};
  $opt = shift if ref $_[0] eq 'HASH';
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
  system($cmd);
  my $returnCode = ${^CHILD_ERROR_NATIVE} / 256;# https://stackoverflow.com/questions/3736320/executing-shell-script-with-system-returns-256-what-does-that-mean
  my ($ansi, $prefix) = $returnCode == 0 ? ('OK') x 2 : ('Err', 'ERR');
  print ansi $ansi, "$prefix: <ansiCmd>$cmd\n" if $opt->{v} && ( 
    $opt->{v} =~ /^all/ ||
    $opt->{v} eq 'ok' && $returnCode == 0 ||
    $opt->{v} eq 'err' && $returnCode != 0 ||
  0);
  return $returnCode;
}

sub gitIsChangedFile($;$) {
  my $fileName = shift || die;
  my $dir = shift;
  my $command;
  if ($dir) {
    $command = 'cd "$dir" && '
  }
  $command .= "git update-index -q --refresh && git diff-index --name-only HEAD -- | grep -Fx \"$fileName\" >/dev/null 2>&1";
  qx/$command/;
  my $result = ${^CHILD_ERROR_NATIVE} / 256; # https://stackoverflow.com/questions/3736320/executing-shell-script-with-system-returns-256-what-does-that-mean
  !$result;
}

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
# =============================================================================

1;