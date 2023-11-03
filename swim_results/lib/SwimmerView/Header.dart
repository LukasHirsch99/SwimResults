import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/model/Swimmer.dart';

class SwimmerPageHeader extends StatefulWidget {
  final Swimmer swimmer;

  const SwimmerPageHeader(this.swimmer, {super.key});

  @override
  State<SwimmerPageHeader> createState() => _SwimmerPageHeaderState();
}

class _SwimmerPageHeaderState extends State<SwimmerPageHeader>
    with TickerProviderStateMixin {
  late TabController tabController;

  @override
  void initState() {
    super.initState();
    tabController = TabController(length: 2, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: MediaQuery.of(context).size.width,
      decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [colorScheme.primary, secondary],
          ),
          boxShadow: const [
            BoxShadow(
              color: Color.fromARGB(70, 0, 0, 0),
              blurRadius: 10,
              offset: Offset(0, 6),
            )
          ]),
      child: Column(
        children: [
          Container(
            margin: const EdgeInsets.only(top: 20),
            alignment: Alignment.centerLeft,
            child: IconButton(
              onPressed: () => Navigator.pop(context),
              icon: const Icon(
                Icons.menu_rounded,
                size: 48,
                color: Colors.white,
              ),
            ),
          ),
          Container(
            margin: const EdgeInsets.only(top: 10),
            child: Text(
              widget.swimmer.name,
              style: const TextStyle(
                  color: Colors.white,
                  fontSize: 30,
                  fontWeight: FontWeight.w600),
            ),
          ),
          Container(
            margin: const EdgeInsets.only(top: 10),
            child: Text(
              widget.swimmer.club!.name,
              style: const TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.w500),
            ),
          ),
          Container(
            margin: const EdgeInsets.only(top: 10, bottom: 10),
            child: Text(
              "${widget.swimmer.birthyear} - ${widget.swimmer.gender!.capitalize()}",
              style: const TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.w500),
            ),
          ),
          TabBar(
            controller: tabController,
            indicatorSize: TabBarIndicatorSize.label,
            unselectedLabelColor: Colors.white,
            indicator: const UnderlineTabIndicator(
                insets: EdgeInsets.symmetric(horizontal: 8),
                borderSide: BorderSide(width: 4, color: Colors.white)),
            tabs: const [
              Tab(
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      "Starts",
                      style: TextStyle(fontSize: 30),
                    ),
                    SizedBox(
                      width: 10,
                    ),
                    Icon(
                      Icons.pool,
                      size: 30,
                    ),
                  ],
                ),
              ),
              Tab(
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      "Results",
                      style: TextStyle(fontSize: 30),
                    ),
                    SizedBox(
                      width: 10,
                    ),
                    Icon(
                      Icons.sports_score,
                      size: 30,
                    ),
                  ],
                ),
              ),
            ],
          )
        ],
      ),
    );
  }
}
