import 'package:flutter/material.dart';
import 'package:swim_results/components/drawer.dart';
import '../components/globals.dart' as globals;

class ProfileSettingPage extends StatefulWidget implements globals.RouteBase {
  static const String routeName = "/profile";

  const ProfileSettingPage({super.key});

  @override
  _ProfileSettingPageState createState() => _ProfileSettingPageState();

  @override
  String getRouteName() {
    return routeName;
  }
}

class _ProfileSettingPageState extends State<ProfileSettingPage> {
  var firstName = TextEditingController();
  var lastName = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawer: const MyDrawer(),
      appBar: AppBar(
        title: const Text(
          'Profile',
          style: TextStyle(color: Colors.tealAccent),
        ),
        centerTitle: true,
      ),
      body: Container(
          alignment: Alignment.center,
          child: Column(children: <Widget>[
            const SizedBox(height: 100),
            Container(
              decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(100),
                  border: Border.all(color: Colors.tealAccent, width: 7)),
              child: const Icon(Icons.person, color: Colors.tealAccent, size: 140),
            ),
            const SizedBox(height: 50),
            const Divider(
              indent: 50,
              endIndent: 50,
              color: Colors.white,
            ),
            const SizedBox(height: 50),
            Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(
                    color: Colors.tealAccent[400],
                    borderRadius: BorderRadius.circular(5)),
                child: Text(globals.myProfile!.name.capitalize(),
                    style: const TextStyle(
                        fontSize: 30,
                        color: Colors.white,
                        fontWeight: FontWeight.bold))),
            const SizedBox(height: 50),
            Container(
                height: 50,
                decoration:
                    BoxDecoration(borderRadius: BorderRadius.circular(10)),
                child: TextButton(
                  // splashColor: Colors.white.withOpacity(.5),
                  // color: Colors.tealAccent[400],
                  child: const Text(
                    'Edit',
                    style: TextStyle(fontSize: 17),
                  ),
                  onPressed: () {
                    Navigator.pushNamed(
                        context, globals.Routes.loginPage.getRouteName());
                  },
                )),
          ])),
    );
  }
}
