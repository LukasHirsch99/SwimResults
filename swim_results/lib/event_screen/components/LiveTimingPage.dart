import 'package:flutter/material.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/drawer.dart';
import 'package:swim_results/model/LiveTiming.dart';
import 'package:rive/rive.dart' as rive;

class LiveTimingPage extends StatefulWidget {
  final int meetId;
  const LiveTimingPage(this.meetId, {super.key});

  @override
  State<LiveTimingPage> createState() => _LiveTimingPageState();
}

class _LiveTimingPageState extends State<LiveTimingPage> {
  late LiveTiming? liveTiming;

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
        future: fetchData(),
        builder: (context, snapshot) {
          if (!snapshot.hasData || snapshot.data == false) {
            return const Center(
              child: SizedBox(
                width: 200,
                height: 200,
                child: rive.RiveAnimation.asset("assets/loading.riv"),
              ),
            );
          }

          return Scaffold(
            drawer: const MyDrawer(),
            appBar: AppBar(
              leading: IconButton(
                icon: const Icon(Icons.menu),
                onPressed: () => Scaffold.of(context).openDrawer(),
              ),
              title: const Text("Live"),
            ),
            body: ListView.builder(
              itemCount: liveTiming!.lanes.length,
              itemBuilder: (ctx, idx) => liveTiming!.lanes[idx],
            ),
          );
        });
  }

  Future<bool> fetchData() async {
    liveTiming = await SwimResultsApi.getLiveTiming(widget.meetId);
    return true;
  }
}
