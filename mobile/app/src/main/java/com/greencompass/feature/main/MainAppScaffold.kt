package com.greencompass.feature.main

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.compose.currentBackStackEntryAsState
import com.greencompass.core.ui.*
import com.greencompass.navigation.AppRoute

data class BottomNavItem(val route: AppRoute, val label: String, val selectedIcon: ImageVector, val unselectedIcon: ImageVector)

val bottomNavItems = listOf(
    BottomNavItem(AppRoute.Today, "Today", Icons.Filled.WbSunny, Icons.Outlined.WbSunny),
    BottomNavItem(AppRoute.Explore, "Explore", Icons.Filled.Explore, Icons.Outlined.Explore),
    BottomNavItem(AppRoute.Updates, "Updates", Icons.Filled.Notifications, Icons.Outlined.Notifications),
    BottomNavItem(AppRoute.Report, "Report", Icons.Filled.EditNote, Icons.Outlined.EditNote),
    BottomNavItem(AppRoute.Profile, "Profile", Icons.Filled.Person, Icons.Outlined.Person)
)

@Composable
fun MainAppScaffold(
    navController: NavController,
    content: @Composable () -> Unit
) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentRouteName = navBackStackEntry?.destination?.route?.substringBefore("?") ?: ""
    
    // Only show bottom nav if we are on one of the main tab routes
    val showBottomNav = bottomNavItems.any { it.route::class.simpleName == currentRouteName }

    Scaffold(
        bottomBar = {
            if (showBottomNav) {
                Column {
                    // Top border: Stone (as per brief)
                    androidx.compose.material3.Divider(color = GreenCompassColors.Stone, thickness = 1.dp)
                    NavigationBar(
                        containerColor = Color.White,
                        tonalElevation = 0.dp,
                        contentColor = GreenCompassColors.Charcoal
                    ) {
                        bottomNavItems.forEach { item ->
                            val selected = currentRouteName == item.route::class.simpleName
                            NavigationBarItem(
                                icon = {
                                    Icon(
                                        imageVector = if (selected) item.selectedIcon else item.unselectedIcon,
                                        contentDescription = item.label
                                    )
                                },
                                label = { Text(text = item.label, style = GreenCompassTypography.labelMedium) },
                                selected = selected,
                                onClick = {
                                    navController.navigate(item.route) {
                                        popUpTo(navController.graph.findStartDestination().id) {
                                            saveState = true
                                        }
                                        launchSingleTop = true
                                        restoreState = true
                                    }
                                },
                                colors = NavigationBarItemDefaults.colors(
                                    selectedIconColor = GreenCompassColors.ForestGreen,
                                    selectedTextColor = GreenCompassColors.ForestGreen,
                                    unselectedIconColor = GreenCompassColors.MutedText,
                                    unselectedTextColor = GreenCompassColors.MutedText,
                                    indicatorColor = Color.Transparent
                                )
                            )
                        }
                    }
                }
            }
        },
        containerColor = GreenCompassColors.WarmWhite
    ) { paddingValues ->
        Box(modifier = Modifier.padding(paddingValues)) {
            content()
        }
    }
}
