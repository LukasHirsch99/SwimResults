import 'package:flutter/material.dart';

// const colorScheme = darkColorScheme;

const textColor = Color(0xFFe7f3f7);
const background = Color(0xFF081519);
const secondaryBackground = Color(0xFF141E21);
const primaryColor = Color(0xFF97cee0);
const primaryFgColor = Color(0xFF081519);
const secondaryColor = Color(0xFF246f89);
const secondaryFgColor = Color(0xFFe7f3f7);
const accentColor = Color(0xFF43aed3);
const accentFgColor = Color(0xFF081519);
  
const colorScheme = ColorScheme(
  brightness: Brightness.dark,
  primary: primaryColor,
  onPrimary: primaryFgColor,
  secondary: secondaryColor,
  onSecondary: secondaryFgColor,
  tertiary: accentColor,
  onTertiary: accentFgColor,
  surface: background,
  onSurface: textColor,
  surfaceTint: secondaryBackground,
  error: Brightness.dark == Brightness.light ? Color(0xffB3261E) : Color(0xffF2B8B5),
  onError: Brightness.dark == Brightness.light ? Color(0xffFFFFFF) : Color(0xff601410),
);
