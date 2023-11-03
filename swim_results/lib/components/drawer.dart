import 'package:flutter/material.dart';
import 'globals.dart';

class MyDrawer extends StatelessWidget {
  const MyDrawer({super.key});

  @override
  Widget build(BuildContext context) {
    return Drawer(
      child: ListView(
        children: [
          GestureDetector(
            child: UserAccountsDrawerHeader(
              accountName: Text(myProfile!.name),
              currentAccountPicture: const CircleAvatar(
                  child: Icon(
                Icons.person,
                size: 40,
                color: Colors.tealAccent,
              )),
              accountEmail: null,
            ),
            onTap: () => Navigator.pushReplacementNamed(
                context, Routes.profilePage.getRouteName()),
          ),
          _createDrawerItem(
              Icons.home,
              "Events",
              () => Navigator.pushReplacementNamed(
                  context, Routes.homepage.getRouteName())),
          _createDrawerItem(
              Icons.watch,
              "Records",
              () => Navigator.pushReplacementNamed(
                  context, Routes.recordPage.getRouteName())),
        ],
      ),
    );
  }

  Widget _createDrawerItem(
      IconData icon, String text, GestureTapCallback onTap) {
    return ListTile(
      title: Row(
        children: <Widget>[
          Icon(icon),
          Padding(
            padding: const EdgeInsets.only(left: 8.0),
            child: Text(text),
          )
        ],
      ),
      onTap: onTap,
    );
  }
}
