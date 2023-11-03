import 'package:flutter/material.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/drawer.dart';
import '../components/globals.dart' as globals;

// ignore: must_be_immutable
class RecordPage extends StatelessWidget implements globals.RouteBase {
  late Map records;
  static const String routeName = '/records';

  RecordPage({super.key});
  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      future: getRecords(),
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          if (snapshot.data == true) {
            return Scaffold(
                drawer: const MyDrawer(),
                appBar: AppBar(
                  title: const Text(
                    'Records',
                    style: TextStyle(color: Colors.tealAccent),
                  ),
                  centerTitle: true,
                ),
                body: SingleChildScrollView(
                  child: Column(
                    children: [
                      ExpansionTile(
                          initiallyExpanded: true,
                          title: const Text('25m'),
                          children: [
                            for (var r in records['25m'])
                              Column(
                                children: [
                                  const Divider(),
                                  ListTile(
                                    title: Text(r['discipline'],
                                        style: const TextStyle(
                                            fontWeight: FontWeight.bold)),
                                    leading: Padding(
                                      padding:
                                          const EdgeInsets.symmetric(vertical: 16),
                                      child: Text(r['time'],
                                          style: const TextStyle(
                                              fontWeight: FontWeight.bold,
                                              color: Colors.tealAccent)),
                                    ),
                                    trailing: Text(r['date'],
                                        style: const TextStyle(
                                            fontWeight: FontWeight.bold)),
                                    subtitle: Text(
                                      r['location'],
                                      style: const TextStyle(color: Colors.white),
                                    ),
                                  ),
                                ],
                              )
                          ]),
                      ExpansionTile(
                          initiallyExpanded: true,
                          title: const Text('50m'),
                          children: [
                            for (var r in records['50m'])
                              Column(
                                children: [
                                  const Divider(),
                                  ListTile(
                                    title: Text(r['discipline'],
                                        style: const TextStyle(
                                            fontWeight: FontWeight.bold)),
                                    leading: Padding(
                                      padding:
                                          const EdgeInsets.symmetric(vertical: 16),
                                      child: Text(r['time'],
                                          style: const TextStyle(
                                              fontWeight: FontWeight.bold,
                                              color: Colors.tealAccent)),
                                    ),
                                    trailing: Text(r['date'],
                                        style: const TextStyle(
                                            fontWeight: FontWeight.bold)),
                                    subtitle: Text(
                                      r['location'],
                                      style: const TextStyle(color: Colors.white),
                                    ),
                                  ),
                                ],
                              )
                          ]),
                    ],
                  ),
                ));
          }
        }
        return const Scaffold(body: Center(child: CircularProgressIndicator()));
      },
    );
  }

  getRecords() async {
    records = await SwimResultsApi.getRecords();
    return records.isNotEmpty;
  }

  @override
  String getRouteName() => routeName;
}
