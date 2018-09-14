package BwCore;
use v5.18;
use strict;
use warnings;

my $selfFileSpec;
BEGIN {
  use File::Find qw/find/;
  my $carpAlwaysIsInstalled;
  no warnings 'File::Find';
  find { wanted => sub { $carpAlwaysIsInstalled = 1 if /\/Carp\/Always(?:\.pm)?$/ }, no_chdir => 1 }, @INC;
  if ( $carpAlwaysIsInstalled ) {
    require Carp::Always;# https://metacpan.org/pod/Carp::Always
    Carp::Always->import;
  }

  use Exporter;
  use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
  $VERSION = 1.00;
  @ISA = qw(Exporter);
  @EXPORT = ();
  @EXPORT_OK = ();
  %EXPORT_TAGS = ();

  my @caller = caller(0);
  $selfFileSpec = $caller[1];
  open(IN, $selfFileSpec);
  while (<IN>) {
    push @EXPORT, $1 if /^\s*sub\s+([a-z][\w\d]*)/;
  }
  close(IN);

  use Hash::Ordered; # https://metacpan.org/pod/Hash::Ordered

  use Data::Dumper;
  push @EXPORT, qw/Dumper/;

  use BwAnsi;
  push @EXPORT, qw/ansi/;

  use BwParams;
  push @EXPORT, qw/run preprocessDefs processParams/;
}

# =============================================================================

sub docker {
  # TODO: bw_install docker --silentIfAlreadyInstalled || return $?
  unshift @_, 'docker';
  if ( $ENV{OSTYPE} =~ m/^linux/ ) {
    unshift @_, 'sudo';
  } elsif ( $ENV{OSTYPE} !~ m/^darwin/ ) {
    die ansi 'Err', "ERR: Неожиданный тип OS <ansiPrimaryLiteral>$ENV{OSTYPE}";
  }
  !&execCmd;
}

sub execCmd {
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
  print ansi "<ansiCmd>$cmd<ansi> . . .\n";
  system($cmd);
  my $returnCode = ${^CHILD_ERROR_NATIVE} / 256;# https://stackoverflow.com/questions/3736320/executing-shell-script-with-system-returns-256-what-does-that-mean
  my ($ansi, $prefix) = $returnCode == 0 ? ('OK') x 2 : ('Err', 'ERR');
  print ansi $ansi, "$prefix: <ansiCmd>$cmd\n";
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