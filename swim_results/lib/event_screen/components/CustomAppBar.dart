import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';

class CustomAppBar extends StatelessWidget {
  final List<Widget> children;
  final List<Widget> tabs;
  const CustomAppBar({super.key, required this.children, required this.tabs});

  @override
  Widget build(BuildContext context) {
    return SliverAppBar(
      floating: true,
      snap: true,
      elevation: 0,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.all(0),
        centerTitle: true,
        title: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: children,
        ),
        background: Container(
          decoration: BoxDecoration(
            color: Colors.white,
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [colorScheme.primary, colorScheme.background],
            ),
          ),
        ),
      ),
      bottom: TabBar(
        indicatorSize: TabBarIndicatorSize.label,
        labelColor: colorScheme.tertiary,
        indicator: UnderlineTabIndicator(
          insets: const EdgeInsets.symmetric(horizontal: 8),
          borderSide: BorderSide(
            width: 4,
            color: colorScheme.tertiary,
          ),
        ),
        tabs: tabs,
      ),
    );
  }
}
