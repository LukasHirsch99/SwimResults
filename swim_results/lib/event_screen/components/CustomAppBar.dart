import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';

class CustomAppBar extends StatelessWidget {
  final Widget title;
  final List<Tab> tabs;
  final double prefferedHeight; 
  final TabController controller;
  const CustomAppBar({
    super.key,
    required this.title,
    required this.tabs,
    required this.controller,
    this.prefferedHeight = 50,
  });

  @override
  Widget build(BuildContext context) {
    return SliverAppBar(
      automaticallyImplyLeading: false,
      floating: true,
      snap: true,
      elevation: 0,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.all(0),
        centerTitle: true,
        title: Center(child: title),
        background: Container(
          decoration: BoxDecoration(
            borderRadius:
                const BorderRadius.vertical(bottom: Radius.circular(10)),
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [colorScheme.secondary, colorScheme.surfaceTint],
            ),
          ),
        ),
      ),
      bottom: PreferredSize(
        preferredSize: Size(0, prefferedHeight),
        child: TabBar(
          controller: controller,
          indicatorSize: TabBarIndicatorSize.tab,
          labelColor: colorScheme.onSurface,
          dividerHeight: 0,
          indicator: RoundEdgeIndicator(
            color: colorScheme.secondary,
            radius: 4,
          ),
          tabs: tabs,
        ),
      ),
    );
  }
}

class RoundEdgeIndicator extends Decoration {
  final BoxPainter _painter;

  RoundEdgeIndicator({required Color color, required double radius})
      : _painter = _CirclePainter(color, radius);

  @override
  BoxPainter createBoxPainter([VoidCallback? onChanged]) => _painter;
}

class _CirclePainter extends BoxPainter {
  final Paint _paint;
  final double radius;

  _CirclePainter(Color color, this.radius)
      : _paint = Paint()
          ..color = color
          ..isAntiAlias = true;

  @override
  void paint(Canvas canvas, Offset offset, ImageConfiguration cfg) {
    canvas.drawRRect(
        RRect.fromRectAndCorners(
          Rect.fromLTWH(
            offset.dx,
            0,
            cfg.size!.width,
            cfg.size!.height,
          ),
          bottomLeft: Radius.lerp(const Radius.circular(10),
              const Radius.circular(0), offset.dx / cfg.size!.width)!,
          bottomRight: Radius.lerp(const Radius.circular(0),
              const Radius.circular(10), offset.dx / cfg.size!.width)!,
        ),
        _paint);
  }
}
