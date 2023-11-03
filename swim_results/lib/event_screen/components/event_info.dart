import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:url_launcher/url_launcher.dart';

class MeetOverview extends StatelessWidget implements RouteBase {
  const MeetOverview(this.m, {super.key});
  static const String routeName = '/meetOverview';
  final Meet m;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            Image.network(m.image!),
            const SizedBox(height: 40),
            infoItem(
              context: context,
              child: Text(
                m.name,
                style: const TextStyle(
                    color: Colors.white,
                    fontSize: 20,
                    fontWeight: FontWeight.w300),
                textAlign: TextAlign.center,
              ),
            ),
            GestureDetector(
              onTap: _launchMap,
              child: infoItem(
                context: context,
                child: const Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    Text(
                      "Maps-Link",
                      style: TextStyle(
                          fontSize: 20,
                          color: Colors.white,
                          fontWeight: FontWeight.w300),
                      textAlign: TextAlign.center,
                    ),
                    Icon(
                      Icons.map,
                      color: Colors.white,
                    )
                  ],
                ),
              ),
            ),
            infoItem(
              context: context,
              child: Column(
                children: [
                  const Text('Invitations',
                      style: TextStyle(color: Colors.white, fontSize: 20)),
                  _invitations(),
                ],
              ),
            ),
            infoItem(
              context: context,
              child: Column(
                children: [
                  const Text(
                    'Registration Deadline',
                    style: TextStyle(color: Colors.white, fontSize: 20),
                  ),
                  Text(
                    DateFormat("dd.MM.yy HH:mm").format(m.deadline),
                    style: const TextStyle(fontSize: 20, color: Colors.white),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  infoItem({required BuildContext context, required Widget child}) {
    return Container(
        width: MediaQuery.of(context).size.width,
        margin: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10),
          color: colorScheme.primary,
          boxShadow: const [
            BoxShadow(
              color: Colors.grey,
              blurRadius: 5,
              offset: Offset(0, 4),
            ),
          ],
        ),
        child: child);
  }

  _launchMap() async {
    final url = Uri.parse('https://www.google.com/maps/search/${m.address}');

    if (await canLaunchUrl(url)) {
      await launchUrl(url);
    } else {
      throw 'Could not launch $url';
    }
  }

  _invitations() {
    List<Widget> children = [];
    for (String c in m.invitations) {
      children.add(IconButton(
          onPressed: () {
            launchUrl(Uri.parse(c));
          },
          icon: const Icon(
            Icons.picture_as_pdf,
            color: Colors.white,
          )));
    }
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: children,
    );
  }

  @override
  String getRouteName() {
    return routeName;
  }
}
