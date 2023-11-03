import 'package:flutter/material.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/event_screen/components/SwimmerPage.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/model/Swimmer.dart';

class SwimmerSearchPage extends StatefulWidget {
  final Meet meet;
  const SwimmerSearchPage(this.meet, {super.key});

  @override
  State<SwimmerSearchPage> createState() => _SwimmerSearchPageState();
}

class _SwimmerSearchPageState extends State<SwimmerSearchPage> {
  SearchController swimmerController = SearchController();
  String currentSearchText = "";
  List<Swimmer> lastSearchResults = [];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Swimmer Search"), automaticallyImplyLeading: false, centerTitle: true,),
      body: Column(children: [
        SearchAnchor(
          searchController: swimmerController,
          builder: (BuildContext context, SearchController controller) {
            return SearchBar(
              controller: swimmerController,
              onTap: () => swimmerController.openView(),
              onChanged: (_) => swimmerController.openView(),
              leading: const Icon(Icons.search),
              hintText: "Name",
            );
            // return IconButton(
            //   onPressed: () => controller.openView(),
            //   icon: const Icon(Icons.search),
            // );
          },
          suggestionsBuilder: (BuildContext context, SearchController controller) async {
            if (currentSearchText != swimmerController.text) {
              currentSearchText = swimmerController.text;
              lastSearchResults = await SwimResultsApi.getSwimmersByNameForEvent(widget.meet.id, swimmerController.text);
            }

            return [
              for (Swimmer s in lastSearchResults)
                ListTile(
                  title: Text(s.name),
                  onTap: () => setState(
                    () {
                      swimmerController.closeView(s.name);
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => SwimmerView(s, widget.meet),
                        ),
                      );
                    },
                  ),
                )
            ];
          },
        ),
      ]),
    );
  }
}
