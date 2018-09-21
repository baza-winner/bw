package BwAnsi;
use v5.18;
use strict;
use warnings;
use Exporter;
use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
$VERSION = 1.00;
@ISA = qw(Exporter);
@EXPORT_OK = ();
%EXPORT_TAGS = ();
@EXPORT = qw/ ansi /;

# =============================================================================

use Data::Dumper;

# =============================================================================

my %ansi;

sub _ansiHelper($) {
  my $key = shift;
  if (!$ansi{$key}) {
    die "ansi '${key}' not found";
  } elsif (ref $ansi{$key} eq 'CODE') {
    $ansi{$key} = &_ansi($ansi{$key}->());
  } else {
    $ansi{$key};
  }
}

sub _ansi($) {
  join '', $ansi{Reset},
  map { _ansiHelper $_ }
  grep { $_ }
  split /[^\w\d]+/, shift
}

sub ansi($;$) {
  my $ansiDefault = 2 == scalar @_ ? _ansi shift : $ansi{Reset};
  $_ = shift || '';
  no strict 'refs';
  s/<ansi([^>]*)>/ $1 ? _ansiHelper $1 : ( $ansiDefault || die "can not use 'ansi' for ansiDefault is not set" )/ge;
  s/^(\s*)/$1$ansiDefault/;
  s/(\s*)$/$ansi{Reset}$1/;
  $_;
}

# https://www.perlmonks.org/?node_id=509827
%ansi = (
  Reset => "\e[0m", # https://superuser.com/questions/33914/why-doesnt-echo-support-e-escape-when-using-the-e-argument-in-macosx/33950#33950
  Bold => "\e[1m",
  Dim => "\e[2m",
  Italic => "\e[3m",
  Underline => "\e[4m",
  Blink => "\e[5m",
#   6 => "\e[6m",
  Invert => "\e[7m",
  Hidden => "\e[8m",
  Strike => "\e[9m",

  ResetBold => "\e[22m", # именно так, а не 20
  ResetDim => "\e[22m",
  ResetItalic => "\e[23m",
  ResetUnderline => "\e[24m",
  ResetBlink => "\e[25m",
#   Reset6 => "\e[26m",
  ResetInvert => "\e[27m",
  ResetHidden => "\e[28m",
  ResetStrike => "\e[29m",

  Black => "\e[30m",
  Red => "\e[31m",
  Green => "\e[32m",
  Yellow => "\e[33m",
  Blue => "\e[34m",
  Magenta => (
    sub {
      my $colors=`tput colors 2>/dev/null`;
      # https://stackoverflow.com/a/32149079
      !${^CHILD_ERROR_NATIVE} && $colors >= 256 ?
        "\e[38;5;201m" : # https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences/33206814#33206814
        "\e[35m"
      ;
    }
  )->(),
  Cyan => "\e[36m",
  LightGray => "\e[37m",
  LightGrey => sub { "LightGray" },
  Default => "\e[39m",
  DarkGray => "\e[90m",
  DarkGrey => sub { "DarkGray" },
  LightRed => "\e[91m",
  LightGreen => "\e[92m",
  LightYellow => "\e[93m",
  LightBlue => "\e[94m",
  LightMagenta => "\e[95m",
  LightCyan => "\e[96m",
  White => "\e[97m",
  Header => sub { "Bold;LightGray" },
  Url => sub { "Blue;Underline" },
  Cmd => sub { "White;Bold" },
  FileSpec => sub { "White;Bold" },
  Dir => sub { "White;Bold" },
  Err => sub { "Red;Bold" },
  Warn => sub { "Yellow;Bold" },
  Will => sub { "Yellow;Bold" },
  OK => sub { "Green;Bold" },
  Outline => sub { "Magenta;Bold" }, #   Outline => sub { "'ansiMagenta''ansiResetBold'"
  Debug => sub { "Blue;ResetBold" },
  PrimaryLiteral => sub { "Cyan;Bold" },
  SecondaryLiteral => sub { "Cyan;ResetBold" },
);

foreach my $key (keys %ansi) {
  if (ref $ansi{$key} eq 'CODE') {
    $ansi{$key} = _ansi $ansi{$key}();
  }
}

1;
