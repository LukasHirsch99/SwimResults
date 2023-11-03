import 'package:flutter/material.dart';
import 'package:swim_results/event_screen/components/personal_result_item.dart';
import 'package:swim_results/event_screen/components/personal_start_item.dart';
import 'package:swim_results/model/Meet.dart';

class SavedSwimmerInfo extends StatelessWidget {
  final List<dynamic> savedSwimmers;
  final Meet m;
  const SavedSwimmerInfo(this.savedSwimmers, this.m, {super.key});

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: savedSwimmers.length,
      itemBuilder: (cntx, i) {
        return ExpansionTile(title: Text(savedSwimmers[i]['name']), children: [
          ExpansionTile(
            title: const Text('Starts'),
            children: [
              for (var s in savedSwimmers[i]['stats']['starts'])
                PersonalStartItem(
                  m.id.toString(),
                  s,
                )
            ],
          ),
          ExpansionTile(
            title: const Text('Results'),
            children: [
              for (var r in savedSwimmers[i]['stats']['results'])
                PersonalResultItem(
                  m.id.toString(),
                  r,
                )
            ],
          ),
        ]);
      },
    );
  }
}
