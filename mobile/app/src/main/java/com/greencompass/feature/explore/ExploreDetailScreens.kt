package com.greencompass.feature.explore

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DetailScaffold(title: String, onBack: () -> Unit, content: @Composable (PaddingValues) -> Unit) {
    GreenCompassScaffold(
        title = title,
        navigationIcon = {
            IconButton(onClick = onBack) {
                androidx.compose.material3.Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        content(paddingValues)
    }
}

@Composable
fun SectionBlock(title: String, content: String, modifier: Modifier = Modifier) {
    Column(modifier = modifier.padding(vertical = AppSpacing.md)) {
        Text(text = title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
        Spacer(modifier = Modifier.height(AppSpacing.xs))
        Text(text = content, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
    }
    Divider(color = GreenCompassColors.Stone)
}

// Screen 26: Local Outlook
@Composable
fun LocalOutlookScreen(onBack: () -> Unit) {
    DetailScaffold(title = "Local outlook", onBack = onBack) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            item {
                Text(text = "Lower Valley", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Conditions may change later today.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                SectionBlock(title = "Today", content = "Partly cloudy\nConditions may become unsettled later.")
                SectionBlock(title = "Tomorrow", content = "Mostly cloudy\nConditions may remain unsettled.")
                
                Text(text = "Updated 10:00", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(top = AppSpacing.md))
                Text(text = "Information provided by\nLocal weather authority", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}

// Screen 27: Seasonal Information
@Composable
fun SeasonalInformationScreen(onBack: () -> Unit) {
    DetailScaffold(title = "Seasonal information", onBack = onBack) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            item {
                Text(text = "Lower Valley", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Conditions may vary during the coming weeks.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                SectionBlock(title = "Rainfall", content = "Rainfall may be below the usual range.")
                SectionBlock(title = "Temperature", content = "Temperatures may remain higher than usual.")
                SectionBlock(title = "What to watch", content = "Conditions may change as the season develops.")
                
                Text(text = "Updated today", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(top = AppSpacing.md))
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}

// Screen 28: Water Outlook
@Composable
fun WaterOutlookScreen(onBack: () -> Unit) {
    DetailScaffold(title = "Water outlook", onBack = onBack) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            item {
                Text(text = "Lower Valley", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Water availability may decline\nover the next two weeks.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                SectionBlock(title = "Current outlook", content = "Information for nearby water sources.")
                SectionBlock(title = "Recent change", content = "How conditions have changed recently.")
                SectionBlock(title = "Community updates", content = "Reports shared near this place.")
                
                Text(text = "Updated today", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(top = AppSpacing.md, bottom = AppSpacing.md))
                
                // Mandatory subdued disclaimer
                Text(
                    text = "This information describes local water conditions.\nIt does not confirm that water is safe to drink.",
                    style = GreenCompassTypography.bodySmall,
                    color = GreenCompassColors.MutedText,
                    modifier = Modifier.padding(vertical = AppSpacing.md)
                )
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}

// Screen 29: Land and Ecosystems
@Composable
fun LandEcosystemsScreen(onBack: () -> Unit) {
    DetailScaffold(title = "Land and ecosystems", onBack = onBack) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            item {
                Text(text = "Lower Valley", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Environmental conditions have changed\nin parts of this area.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                SectionBlock(title = "Vegetation", content = "Recent conditions compared with the usual range.")
                SectionBlock(title = "Surface conditions", content = "Recent information for the surrounding area.")
                SectionBlock(title = "What to watch", content = "Changes may continue during the coming weeks.")
                
                Text(text = "Updated today", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(top = AppSpacing.md))
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}

// Screen 30: Food and Agriculture
@Composable
fun FoodAgricultureScreen(onBack: () -> Unit) {
    DetailScaffold(title = "Food and agriculture", onBack = onBack) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            item {
                Text(text = "Lower Valley", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Information relevant to the current season.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                SectionBlock(title = "Seasonal timing", content = "Current agricultural periods for this area.")
                SectionBlock(title = "Local conditions", content = "Information that may affect the season.")
                SectionBlock(title = "Guidance", content = "Approved information from local agricultural teams.")
                
                Text(text = "Updated today", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(top = AppSpacing.md))
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}

// Screen 31: Community Updates
@Composable
fun CommunityUpdatesScreen(onBack: () -> Unit) {
    DetailScaffold(title = "Community updates", onBack = onBack) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            item {
                Text(text = "Lower Valley", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Recent information shared by people\nand local teams nearby.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
            }
            
            // Demo Card
            item {
                Card(
                    shape = RoundedCornerShape(16.dp),
                    colors = CardDefaults.cardColors(containerColor = Color.White),
                    border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
                ) {
                    Column(modifier = Modifier.padding(AppSpacing.lg)) {
                        Text(text = "Water point update", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                        Spacer(modifier = Modifier.height(AppSpacing.xs))
                        Text(text = "A local water source was reported\nas unavailable.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
                        Spacer(modifier = Modifier.height(AppSpacing.sm))
                        Text(text = "Under review · Updated 09:30", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
                        Spacer(modifier = Modifier.height(AppSpacing.sm))
                        TextLinkButton(text = "Read update", onClick = { })
                    }
                }
            }
            
            // Empty State (Commented out for demo, but structure is here per brief)
            /*
            item {
                EmptyStateView(
                    icon = Icons.Outlined.Info,
                    title = "No community updates yet",
                    message = "There are no recent updates for this place."
                )
                PrimaryButton(text = "Share an update", onClick = { }, modifier = Modifier.padding(horizontal = AppSpacing.lg))
            }
            */
            
            item {
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}
