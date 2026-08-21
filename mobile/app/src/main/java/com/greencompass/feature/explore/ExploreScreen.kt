package com.greencompass.feature.explore

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.*
import com.greencompass.feature.today.TodayViewModel
import com.greencompass.navigation.AppRoute

data class ExploreTopic(val key: String, val title: String, val description: String, val icon: ImageVector, val route: AppRoute)

@Composable
fun ExploreScreen(
    viewModel: TodayViewModel = hiltViewModel(),
    onNavigate: (AppRoute) -> Unit
) {
    val state by viewModel.uiState.collectAsState()
    
    GreenCompassScaffold(
        title = "Explore",
        showPlaceSwitcher = true,
        placeName = state.data?.placeName ?: "Lower Valley"
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Information relevant to your selected place.",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            val topics = listOf(
                ExploreTopic("local", "Local outlook", "Conditions for the coming days", Icons.Outlined.WbSunny, AppRoute.LocalOutlook),
                ExploreTopic("seasonal", "Seasonal information", "Changes that may affect your area", Icons.Outlined.CalendarToday, AppRoute.SeasonalInformation),
                ExploreTopic("water", "Water outlook", "Information about nearby water conditions", Icons.Outlined.WaterDrop, AppRoute.WaterOutlook),
                ExploreTopic("land", "Land and ecosystems", "Changes in the surrounding environment", Icons.Outlined.Landscape, AppRoute.LandEcosystems),
                ExploreTopic("food", "Food and agriculture", "Information relevant to the current season", Icons.Outlined.Agriculture, AppRoute.FoodAgriculture),
                ExploreTopic("community", "Community updates", "Information shared by people nearby", Icons.Outlined.People, AppRoute.CommunityUpdates)
            )

            LazyVerticalGrid(
                columns = GridCells.Fixed(2),
                horizontalArrangement = Arrangement.spacedBy(AppSpacing.sm),
                verticalArrangement = Arrangement.spacedBy(AppSpacing.sm),
                modifier = Modifier.weight(1f)
            ) {
                items(topics) { topic ->
                    ExploreTopicCard(topic = topic, onClick = { onNavigate(topic.route) })
                }
            }
            
            Spacer(modifier = Modifier.height(AppSpacing.xxxl))
        }
    }
}

@Composable
private fun ExploreTopicCard(topic: ExploreTopic, onClick: () -> Unit) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .aspectRatio(1f)
            .clickable { onClick() },
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Column(
            modifier = Modifier.padding(AppSpacing.md).fillMaxSize(),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            Icon(
                imageVector = topic.icon,
                contentDescription = null,
                tint = GreenCompassColors.ForestGreen,
                modifier = Modifier.size(32.dp)
            )
            Column {
                Text(text = topic.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                Spacer(modifier = Modifier.height(AppSpacing.xxs))
                Text(text = topic.description, style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
            }
        }
    }
}
