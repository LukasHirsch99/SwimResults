import 'dart:async';
import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/home_screen/home_screen.dart';
import 'package:swim_results/model/Club.dart';
import 'package:swim_results/model/Swimmer.dart';
import '../components/globals.dart' as globals;

class LoginPage extends StatefulWidget implements globals.RouteBase {
  static const String routeName = '/login';

  const LoginPage({super.key});
  @override
  State<LoginPage> createState() => _LoginPageState();

  @override
  String getRouteName() {
    return routeName;
  }
}

class _LoginPageState extends State<LoginPage> {
  Swimmer? swimmer;
  Club? club;
  SearchController swimmerController = SearchController();
  bool validName = true, loading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: Navigator.canPop(context)
            ? IconButton(
                icon: Icon(Icons.arrow_back_ios, color: colorScheme.primary),
                onPressed: () => Navigator.pop(context),
              )
            : null,
        title: Text('Login',
            style: TextStyle(
              color: colorScheme.primary,
            )),
        centerTitle: true,
      ),
      body: Column(
        children: <Widget>[
          globals.trainerMode ? trainerLogin() : swimmerLogin(),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
            decoration:
                BoxDecoration(borderRadius: BorderRadius.circular(10)),
            child: ElevatedButton(
              style: ButtonStyle(
                backgroundColor: MaterialStateColor.resolveWith(
                    (states) => colorScheme.primary),
              ),
              child: Text(
                globals.trainerMode ? "I'm a Swimmer" : "I'm a Trainer",
                style: const TextStyle(fontSize: 17, color: Colors.white,),
              ),
              onPressed: () => setState(() => globals.trainerMode = !globals.trainerMode)
            ),
          ),
          Container(
            height: 50,
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
            decoration:
                BoxDecoration(borderRadius: BorderRadius.circular(10)),
            child: ElevatedButton(
              style: ButtonStyle(
                backgroundColor: MaterialStateColor.resolveWith(
                    (states) => colorScheme.primary),
              ),
              child: const Text(
                'Login',
                style: TextStyle(fontSize: 17, color: Colors.white),
              ),
              onPressed: () => _login(),
            ),
          ),
          _loginStatus(),
        ],
      ),
    );
  }

  Widget trainerLogin() {
    return Column(
      children: [
        SearchAnchor(
          builder: (BuildContext context, SearchController controller) {
            return SearchBar(
              controller: controller,
              onTap: () => controller.openView(),
              onChanged: (_) => controller.openView(),
              leading: const Icon(Icons.search),
            );
          },
          suggestionsBuilder: (BuildContext context, SearchController controller) async {
            List<Club> clubs = await SwimResultsApi.getClubsByName(controller.text);
            return List.generate(
              clubs.length,
              (i) => ListTile(
                title: Text(clubs[i].name),
                onTap: () => setState(() {
                  controller.closeView(clubs[i].name);
                  club = clubs[i];
                }),
              )
            );
          },
        ),
      ], 
    );
  }

  Widget swimmerLogin() {
    return Column(children: [
      Container(
        margin: const EdgeInsets.symmetric(horizontal: 30),
        child: SearchAnchor(
          searchController: swimmerController,
          builder: (BuildContext context, SearchController controller) {
            return SearchBar(
              controller: swimmerController,
              onTap: () => swimmerController.openView(),
              onChanged: (_) => swimmerController.openView(),
              leading: const Icon(Icons.search),
              hintText: "Name",
            );
          },
          suggestionsBuilder: (BuildContext context, SearchController controller) async {
            List<Swimmer> swimmers = await SwimResultsApi.getSwimmersByName(swimmerController.text);
            return [
              for (Swimmer s in swimmers)
                ListTile(
                  title: Text(s.fullname),
                  onTap: () => setState(() {
                    swimmerController.closeView(s.fullname);
                    swimmer = s;
                    // _login();
                  }),
                )
            ];  
          },
        ),
      ),
    ],);
  }


  @override
  void dispose() {
    swimmerController.dispose();
    super.dispose();
  }

  bool _login() {
    if (!globals.trainerMode) {
      if (swimmer == null) {
        setState(() => validName = false);
        return false;
      }
      globals.setProfile(swimmer!);
      validName = true;
    }
    else {
      if (club == null) {
        setState(() => validName = false);
        return false;
      }
      globals.setClub(club!);
      validName = true;
    }

    if (Navigator.canPop(context)) {
      Navigator.pushAndRemoveUntil(
        context,
        MaterialPageRoute(builder: (cntx) => const HomePage()),
        (Route<dynamic> route) => false
      );
    } else {
      Navigator.push(
        context,
        MaterialPageRoute(builder: (cntx) => const HomePage())
      );
    }
    return true;
  }

  _loginStatus() {
    if (loading) {
      return Container(
        margin: const EdgeInsets.symmetric(vertical: 15),
        child: const CircularProgressIndicator(strokeWidth: 3),
      );
    } else if (validName == false) {
      return Container(
        margin: const EdgeInsets.symmetric(vertical: 15),
        child: const Text(
          'Swimmer not found',
          style: TextStyle(color: Colors.red),
        ),
      );
    } else {
      return const SizedBox();
    }
  }
}
