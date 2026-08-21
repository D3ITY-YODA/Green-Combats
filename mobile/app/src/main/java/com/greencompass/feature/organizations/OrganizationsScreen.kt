package com.greencompass.feature.organizations

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun OrganizationsScreen(
    onBack: () -> Unit,
    onOpenOrganization: () -> Unit,
    onViewRequest: () -> Unit
) {
    GreenCompassScaffold(
        title = "Organizations",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg),
            verticalArrangement = Arrangement.spacedBy(AppSpacing.sm)
        ) {
            item {
                Text(text = "Your organization access", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
            }

            item {
                Card(
                    modifier = Modifier.fillMaxWidth().clickable { onOpenOrganization() },
                    shape = RoundedCornerShape(16.dp),
                    colors = CardDefaults.cardColors(containerColor = Color.White),
                    border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
                ) {
                    Column(modifier = Modifier.padding(AppSpacing.lg)) {
                        Text(text = "Lower Valley Water Authority", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                        Spacer(modifier = Modifier.height(AppSpacing.xxs))
                        Text(text = "Water authority · Water reviewer", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
                        Spacer(modifier = Modifier.height(AppSpacing.sm))
                        TextLinkButton(text = "Open organization", onClick = onOpenOrganization)
                    }
                }
            }

            item {
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(16.dp),
                    colors = CardDefaults.cardColors(containerColor = Color.White),
                    border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
                ) {
                    Column(modifier = Modifier.padding(AppSpacing.lg)) {
                        Text(text = "Eastern Regional Environment Office", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                        Spacer(modifier = Modifier.height(AppSpacing.xxs))
                        Text(text = "Access pending", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.ImportantAmber)
                        Spacer(modifier = Modifier.height(AppSpacing.sm))
                        TextLinkButton(text = "View request", onClick = onViewRequest)
                    }
                }
            }
        }
    }
}
