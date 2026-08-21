package com.greencompass.feature.onboarding

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.ChevronRight
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*
import com.greencompass.domain.model.Organization

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun OrganizationSearchScreen(
    onBack: () -> Unit,
    onSelectOrganization: (Organization) -> Unit,
    onRequestAccess: () -> Unit
) {
    var query by remember { mutableStateOf("") }
    val demoOrgs = listOf(
        Organization("Lower Valley Water Authority", "Water authority", "Lower Valley"),
        Organization("Eastern Regional Environment Office", "Local government", "Eastern Region"),
        Organization("Green Horizons Community Network", "Community organization", "Lower Valley and East Ward")
    )

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Join an organization",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            OutlinedTextField(
                value = query,
                onValueChange = { query = it },
                label = { Text("Search for your organization") },
                leadingIcon = { Icon(Icons.Default.Search, contentDescription = null) },
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.xl),
                shape = RoundedCornerShape(12.dp)
            )

            LazyColumn(verticalArrangement = Arrangement.spacedBy(AppSpacing.sm)) {
                items(demoOrgs) { org ->
                    Surface(
                        modifier = Modifier.fillMaxWidth().clickable { onSelectOrganization(org) },
                        shape = RoundedCornerShape(12.dp),
                        color = GreenCompassColors.WarmWhite,
                        border = BorderStroke(1.dp, GreenCompassColors.Stone)
                    ) {
                        Row(
                            modifier = Modifier.padding(AppSpacing.md).fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Column {
                                Text(text = org.name, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                                Text(text = "${org.type} · ${org.area}", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
                            }
                            Icon(Icons.Default.ChevronRight, contentDescription = null, tint = GreenCompassColors.MutedText)
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = "Can't find it?",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )

            TextLinkButton(
                text = "Request organization access",
                onClick = onRequestAccess
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
