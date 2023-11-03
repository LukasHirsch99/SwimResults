import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/globals.dart' as globals;
import 'package:swim_results/login_screen/login_screen.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:swim_results/home_screen/home_screen.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

Future<void> main() async {
  await Supabase.initialize(
      url: 'https://qeudknoyuvjztxvgbmou.supabase.co',
      anonKey:
          'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU');
  globals.prefs = await SharedPreferences.getInstance();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Swim Results',
      theme: ThemeData(
        colorScheme: colorScheme,
      ), 
      routes: globals.Routes.getRoutes(),
      home: FutureBuilder<bool>(
        future: showProfile(),
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            if (snapshot.data == true) {
              return const HomePage();
            }
            return const LoginPage();
          }
          return const Scaffold(body: Center(child: CircularProgressIndicator()));
        },
      ),
    );
  }

  Future<bool> showProfile() async {
    globals.myProfile = globals.getProfile();

    return (globals.myProfile != null);
  }
}
