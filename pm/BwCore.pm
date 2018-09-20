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
  execCmd
  docker
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

sub execCmd {
  my $opt = ref $_[0] eq 'HASH' ? shift : {};
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
  my $stdout;
  if (hasItem 'stdout', @{_getAsArrayRef($opt->{return})}) {
    $stdout = qx/$cmd/;
  } else {
    system($cmd);
  }
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
      }
    }
    return $errorCode;
  } else {
    my @result = ();
    foreach my $item (@optReturn) {
      if (defined $item) {
        if ($item eq 'stdout') {
          push @result, $stdout;
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
    my ($commentPrefix, $commentSuffix, $commentPreLine, $commentPostLine) = ('') x 4;
    if ($fileSpec =~ m/\.html$/) {
      ($commentPreLine, $commentPostLine) = ( '<!--', '-->' );
    } elsif ($fileSpec =~ m/\.conf$/) {
      $commentPrefix = ('# ');
    } else {
      die Dumper({fileSpec => $fileSpec});
    }

    unless ( exists $opt->{noNotice} or $opt->{noNotice} ) {
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