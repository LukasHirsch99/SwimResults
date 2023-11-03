import 'package:flutter/material.dart';

const colorScheme = darkColorScheme;

const darkColorScheme = ColorScheme(
  brightness: Brightness.dark,
  background: Color(0xFF030c0a),
  onBackground: Color(0xFFe3f7f4),
  primary: Color(0xFF43cbb5),
  onPrimary: Color(0xFF030c0a),
  secondary: Color(0xFF13443c),
  onSecondary: Color(0xFFe3f7f4),
  tertiary: Color(0xFF53d0bb),
  onTertiary: Color(0xFF030c0a),
  surface: Color(0xFF030c0a),
  onSurface: Color(0xFFe3f7f4),
  error: Brightness.dark == Brightness.light
      ? Color(0xffB3261E)
      : Color(0xffF2B8B5),
  onError: Brightness.dark == Brightness.light
      ? Color(0xffFFFFFF)
      : Color(0xff601410),
);

const lightColorScheme = ColorScheme(
  brightness: Brightness.light,
  background: Color(0xFFf3fcfa),
  onBackground: Color(0xFF081c19),
  primary: Color(0xFF34bca5),
  onPrimary: Color(0xFF081c19),
  secondary: Color(0xFFbbece4),
  onSecondary: Color(0xFF081c19),
  tertiary: Color(0xFF2fac97),
  onTertiary: Color(0xFF081c19),
  surface: Color(0xFFf3fcfa),
  onSurface: Color(0xFF081c19),
  error: Brightness.light == Brightness.light
      ? Color(0xffB3261E)
      : Color(0xffF2B8B5),
  onError: Brightness.light == Brightness.light
      ? Color(0xffFFFFFF)
      : Color(0xff601410),
);
